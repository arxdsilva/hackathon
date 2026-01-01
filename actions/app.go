package actions

import (
	"net/http"
	"sync"

	"github.com/arxdsilva/hackathon/locales"
	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/public"
	"github.com/arxdsilva/hackathon/repository"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/gobuffalo/pop/v6"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

type MyApp struct {
	*buffalo.App
}

// Repository returns a repository interface for the given transaction
func (a *MyApp) Repository(tx *pop.Connection) repository.RepositoryInterface {
	return repository.NewRepositoryManager(tx)
}

var (
	myApp   *MyApp
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
		myApp = &MyApp{buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_hackathon_session",
		})}

		// Automatically redirect to SSL
		myApp.Use(myApp.forceSSL())

		// Log request parameters (filters apply).
		myApp.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		myApp.Use(csrf.New)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		myApp.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		myApp.Use(myApp.translations())

		// Load the current user into context and protect routes.
		myApp.Use(myApp.SetCurrentUser)
		myApp.Use(myApp.RequirePasswordReset)
		myApp.Use(myApp.Authorize)

		myApp.GET("/", myApp.HomeHandler)
		myApp.GET("/about", myApp.AboutHandler)
		myApp.GET("/hackathons", myApp.RequireLogin(myApp.HackathonsIndex))
		myApp.GET("/hackathons/new", myApp.RequireRoleOwner(myApp.HackathonsNew))
		myApp.POST("/hackathons", myApp.RequireRoleOwner(myApp.HackathonsCreate))
		myApp.GET("/hackathons/{hackathon_id}", myApp.HackathonsShow)
		myApp.GET("/hackathons/{hackathon_id}/edit", myApp.RequireHackathonOwner(myApp.HackathonsEdit))
		myApp.PUT("/hackathons/{hackathon_id}", myApp.RequireHackathonOwner(myApp.HackathonsUpdate))
		myApp.DELETE("/hackathons/{hackathon_id}", myApp.RequireHackathonOwner(myApp.HackathonsDestroy))
		myApp.GET("/hackathons/{hackathon_id}/projects", myApp.ProjectsIndex)
		myApp.GET("/hackathons/{hackathon_id}/projects/new", myApp.ProjectsNew)
		myApp.POST("/hackathons/{hackathon_id}/projects", myApp.ProjectsCreate)
		myApp.GET("/hackathons/{hackathon_id}/projects/{project_id}", myApp.ProjectsShow)
		myApp.GET("/hackathons/{hackathon_id}/projects/{project_id}/image", myApp.ProjectsImage)
		myApp.PUT("/hackathons/{hackathon_id}/projects/{project_id}/image", myApp.ProjectsUpdateImage)
		myApp.GET("/hackathons/{hackathon_id}/projects/{project_id}/edit", myApp.ProjectsEdit)
		myApp.PUT("/hackathons/{hackathon_id}/projects/{project_id}", myApp.ProjectsUpdate)
		myApp.POST("/hackathons/{hackathon_id}/projects/{project_id}/toggle-presenting", myApp.ProjectsTogglePresenting)
		myApp.POST("/hackathons/{hackathon_id}/projects/{project_id}/join", myApp.ProjectMembershipsCreate)
		myApp.DELETE("/hackathons/{hackathon_id}/projects/{project_id}/leave", myApp.ProjectMembershipsDestroy)
		myApp.GET("/routes", myApp.RequireRoleOwner(myApp.RoutesHandler))
		myApp.GET("/profile", myApp.ProfileShow)
		myApp.GET("/profile/edit", myApp.ProfileEdit)
		myApp.PUT("/profile", myApp.ProfileUpdate)
		myApp.GET("/users/new", myApp.UsersNew)
		myApp.POST("/users", myApp.UsersCreate)
		myApp.GET("/users/{user_id}/edit", myApp.RequireRoleOwner(myApp.UsersEdit)).Name("userEditPath")
		myApp.PUT("/users/{user_id}", myApp.RequireRoleOwner(myApp.UsersUpdate)).Name("userPath")
		myApp.GET("/files", myApp.RequireLogin(myApp.FilesIndex))
		myApp.GET("/files/new", myApp.RequireLogin(myApp.FilesNew))
		myApp.POST("/files", myApp.RequireLogin(myApp.FilesCreate))
		myApp.GET("/files/{file_id}", myApp.RequireLogin(myApp.FilesShow))
		myApp.GET("/files/{file_id}/download", myApp.RequireLogin(myApp.FilesDownload))
		myApp.DELETE("/files/{file_id}", myApp.RequireLogin(myApp.FilesDestroy))
		myApp.GET("/signin", myApp.AuthNew)
		myApp.POST("/signin", myApp.AuthCreate)
		myApp.DELETE("/signout", myApp.AuthDestroy)
		myApp.GET("/reset-password", myApp.ResetPasswordNew)
		myApp.POST("/reset-password", myApp.ResetPasswordCreate)

		// Allow unauthenticated access to Home, About, and Auth endpoints
		myApp.Middleware.Skip(myApp.Authorize, myApp.HomeHandler, myApp.AboutHandler, myApp.UsersNew, myApp.UsersCreate, myApp.AuthNew, myApp.AuthCreate, myApp.ResetPasswordNew, myApp.ResetPasswordCreate)

		// Admin routes
		admin := myApp.Group("/admin")
		admin.Use(myApp.RequireRoleOwner)
		admin.GET("/", myApp.AdminIndex)
		admin.GET("/users", myApp.AdminUsersIndex)
		admin.GET("/users/new", myApp.AdminUsersNew)
		admin.POST("/users", myApp.AdminUsersCreate)
		admin.GET("/users/{user_id}", myApp.AdminUsersShow)
		admin.GET("/users/{user_id}/edit", myApp.AdminUsersEdit)
		admin.PUT("/users/{user_id}", myApp.AdminUsersUpdate)
		admin.DELETE("/users/{user_id}", myApp.AdminUsersDestroy)
		admin.GET("/hackathons", myApp.AdminHackathonsIndex)
		admin.GET("/projects", myApp.AdminProjectsIndex)
		admin.GET("/emails", myApp.AdminEmailsIndex)
		admin.GET("/config", myApp.AdminConfigIndex)
		admin.PUT("/config", myApp.AdminConfigUpdate).Name("adminConfigPath")
		admin.GET("/passwords", myApp.AdminPasswordsIndex)
		admin.GET("/domains", myApp.AdminDomainsIndex)
		admin.POST("/domains", myApp.AdminDomainsCreate)
		admin.PUT("/domains/{domain_id}", myApp.AdminDomainsUpdate)
		admin.DELETE("/domains/{domain_id}", myApp.AdminDomainsDestroy)
		admin.GET("/audit-logs", myApp.AdminAuditLogsIndex)

		myApp.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	})

	return myApp.App
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func (a *MyApp) translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		a.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func (a *MyApp) forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
