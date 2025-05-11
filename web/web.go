package web

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/frontend"
	"github.com/sirupsen/logrus"
	"goji.io/v3"
	"goji.io/v3/pat"
)

var (
	RootMultiplexer      *goji.Mux
	DashboardMultiplexer *goji.Mux

	HTMLTemplates fs.FS = frontend.HTMLTemplates
	StaticFiles   fs.FS = frontend.StaticFiles

	URL string = "https://" + common.ConfigASBWIGHost
)

func Run() {
	initDiscordOauth()
	multiplexer := setupWebRoutes()
	runWebServer(multiplexer)
}

func embedHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := dashboardContextData(w, r)

		tmpl, err := template.ParseFS(frontend.HTMLTemplates, "templates/*.html")
		if err != nil {
			http.Error(w, "Failed to parse templates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.ExecuteTemplate(w, filename, data); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}

func setupWebRoutes() *goji.Mux {
	// start the base routes, such as logins, static files and such
	runRootMultiplexer()

	// Handle dashboard page
	RootMultiplexer.HandleFunc(pat.Get("/dashboard"), embedHTML("dashboard.html"))
	RootMultiplexer.HandleFunc(pat.Get("/dashboard/"), embedHTML("dashboard.html"))

	// Create a sub-mux for dashboard-related routes
	DashboardMultiplexer = goji.SubMux()
	DashboardMultiplexer.Use(validateGuild)

	// Server manage pages
	RootMultiplexer.Handle(pat.New("/dashboard/:server"), DashboardMultiplexer)
	RootMultiplexer.Handle(pat.New("/dashboard/:server/*"), DashboardMultiplexer)

	DashboardMultiplexer.HandleFunc(pat.Get("/manage"), embedHTML("manage.html"))
	DashboardMultiplexer.HandleFunc(pat.Get("/manage/"), embedHTML("manage.html"))

	DashboardMultiplexer.HandleFunc(pat.Get("/manage/core"), embedHTML("core.html"))
	DashboardMultiplexer.HandleFunc(pat.Get("/manage/core/"), embedHTML("core.html"))

	DashboardMultiplexer.HandleFunc(pat.Post("/manage/update-prefix"), handleUpdatePrefix)
	DashboardMultiplexer.HandleFunc(pat.Post("/manage/update-prefix/"), handleUpdatePrefix)

	return RootMultiplexer
}

func runRootMultiplexer() {
	mux := goji.NewMux()
	RootMultiplexer = mux

	mux.Handle(pat.Get("/static/*"), http.FileServer(http.FS(StaticFiles)))

	// Serve the login page
	mux.HandleFunc(pat.Get("/"), embedHTML("index.html"))
	mux.HandleFunc(pat.Get("/login"), handleLogin)
	mux.HandleFunc(pat.Get("/logout"), handleLogout)
	mux.HandleFunc(pat.Get("/confirm"), confirmLogin)

	// Data and service related pages
	mux.HandleFunc(pat.Get("/terms"), embedHTML("terms.html"))
	mux.HandleFunc(pat.Get("/privacy-policy"), embedHTML("privacy.html"))
}

func runWebServer(multiplexer *goji.Mux) {
	logrus.Info("Webserver started on :8085")
	http.ListenAndServe(":8085", multiplexer)
}