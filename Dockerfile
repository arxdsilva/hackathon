# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/

# Go build stage
FROM gobuffalo/buffalo:latest as builder

ENV GOPROXY http://proxy.golang.org

# Install Node.js 18.x (compatible with modern packages)
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install -y nodejs

# Install Go 1.24
RUN wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Install Yarn 1.x (compatible with older Node.js)
RUN curl -o- -L https://yarnpkg.com/install.sh | bash -s -- --version 1.22.19
ENV PATH="$HOME/.yarn/bin:$HOME/.config/yarn/global/node_modules/.bin:$PATH"

RUN mkdir -p /src/hackathon
WORKDIR /src/hackathon

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install dependencies
RUN yarn install

# Build the application
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD exec /bin/app
