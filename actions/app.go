package actions

import (
	"net/http"
	"sync"

	"github.com/arxdsilva/hackathon/locales"
	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_hackathon_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		app.Use(translations())

		// Load the current user into context and protect routes.
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		app.GET("/", HomeHandler)
		app.GET("/about", AboutHandler)
		app.GET("/hackathons", RequireLogin(HackathonsIndex))
		app.GET("/hackathons/new", RequireRoleOwner(HackathonsNew))
		app.POST("/hackathons", RequireRoleOwner(HackathonsCreate))
		app.GET("/hackathons/{hackathon_id}", HackathonsShow)
		app.GET("/hackathons/{hackathon_id}/edit", RequireHackathonOwner(HackathonsEdit))
		app.PUT("/hackathons/{hackathon_id}", RequireHackathonOwner(HackathonsUpdate))
		app.DELETE("/hackathons/{hackathon_id}", RequireHackathonOwner(HackathonsDestroy))
		app.GET("/hackathons/{hackathon_id}/projects", ProjectsIndex)
		app.GET("/hackathons/{hackathon_id}/projects/new", ProjectsNew)
		app.POST("/hackathons/{hackathon_id}/projects", ProjectsCreate)
		app.GET("/hackathons/{hackathon_id}/projects/{project_id}", ProjectsShow)
		app.GET("/hackathons/{hackathon_id}/projects/{project_id}/image", ProjectsImage)
		app.PUT("/hackathons/{hackathon_id}/projects/{project_id}/image", ProjectsUpdateImage)
		app.GET("/hackathons/{hackathon_id}/projects/{project_id}/edit", ProjectsEdit)
		app.PUT("/hackathons/{hackathon_id}/projects/{project_id}", ProjectsUpdate)
		app.POST("/hackathons/{hackathon_id}/projects/{project_id}/join", ProjectMembershipsCreate)
		app.DELETE("/hackathons/{hackathon_id}/projects/{project_id}/leave", ProjectMembershipsDestroy)
		app.GET("/routes", RequireRoleOwner(RoutesHandler))
		app.GET("/profile", ProfileShow)
		app.GET("/profile/edit", ProfileEdit)
		app.PUT("/profile", ProfileUpdate)
		app.GET("/users/new", UsersNew)
		app.POST("/users", UsersCreate)
		app.GET("/users/{user_id}/edit", RequireRoleOwner(UsersEdit)).Name("userEditPath")
		app.PUT("/users/{user_id}", RequireRoleOwner(UsersUpdate)).Name("userPath")
		app.GET("/files", RequireLogin(FilesIndex))
		app.GET("/files/new", RequireLogin(FilesNew))
		app.POST("/files", RequireLogin(FilesCreate))
		app.GET("/files/{file_id}", RequireLogin(FilesShow))
		app.GET("/files/{file_id}/download", RequireLogin(FilesDownload))
		app.DELETE("/files/{file_id}", RequireLogin(FilesDestroy))
		app.GET("/signin", AuthNew)
		app.POST("/signin", AuthCreate)
		app.DELETE("/signout", AuthDestroy)

		// Allow unauthenticated access to Home, About, and Auth endpoints
		app.Middleware.Skip(Authorize, HomeHandler, AboutHandler, UsersNew, UsersCreate, AuthNew, AuthCreate)

		// Admin routes
		admin := app.Group("/admin")
		admin.Use(RequireRoleOwner)
		admin.GET("/", AdminIndex)
		admin.GET("/users", AdminUsersIndex)
		admin.GET("/users/{user_id}", AdminUsersShow)
		admin.GET("/users/{user_id}/edit", AdminUsersEdit)
		admin.PUT("/users/{user_id}", AdminUsersUpdate)
		admin.DELETE("/users/{user_id}", AdminUsersDestroy)
		admin.GET("/hackathons", AdminHackathonsIndex)
		admin.GET("/projects", AdminProjectsIndex)
		admin.GET("/emails", AdminEmailsIndex)
		admin.GET("/config", AdminConfigIndex)
		admin.PUT("/config", AdminConfigUpdate).Name("adminConfigPath")
		admin.GET("/passwords", AdminPasswordsIndex)
		admin.GET("/domains", AdminDomainsIndex)
		admin.POST("/domains", AdminDomainsCreate)
		admin.PUT("/domains/{domain_id}", AdminDomainsUpdate)
		admin.DELETE("/domains/{domain_id}", AdminDomainsDestroy)

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
