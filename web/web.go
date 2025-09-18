package web

import (
	"encoding/json"
	"fmt"
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
	HTMLPages     fs.FS = frontend.HTMLPages
	StaticFiles   fs.FS = frontend.StaticFiles

	URL        string = "https://" + common.ConfigASBWIGHost
	TermsURL   string = common.ConfigTermsURLOverride
	PrivacyURL string = common.ConfigPrivacyURLOverride

	dashboardRoutes []func(*goji.Mux)
)

func Run() {
	initDiscordOauth()
	multiplexer := setupWebRoutes()
	runWebServer(multiplexer)
}

func EmbedHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.New("").Funcs(templateFunctions)

		// Parse templates
		_, err := tmpl.ParseFS(frontend.HTMLTemplates, "templates/*.html")
		if err != nil {
			http.Error(w, "Failed to parse templates", http.StatusInternalServerError)
			return
		}
		// Parse pages
		_, err = tmpl.ParseFS(frontend.HTMLPages, "pages/*.html")
		if err != nil {
			http.Error(w, "Failed to parse pages", http.StatusInternalServerError)
			return
		}

		tmplData, _ := r.Context().Value(CtxKeyTmplData).(TmplContextData)

		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.ExecuteTemplate(w, filename, tmplData); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

func setupWebRoutes() *goji.Mux {
	// start the base routes, such as logins, static files and such
	runRootMultiplexer()

	// Handle dashboard page
	RootMultiplexer.Handle(pat.Get("/dashboard"), userAndManagedGuildsInfoMW(EmbedHTML("dashboard.html")))
	RootMultiplexer.Handle(pat.Get("/dashboard/"), userAndManagedGuildsInfoMW(EmbedHTML("dashboard.html")))

	// Create a sub-mux for dashboard-related routes
	DashboardMultiplexer = goji.SubMux()

	// Middlewares
	DashboardMultiplexer.Use(validateGuild)
	DashboardMultiplexer.Use(userAndManagedGuildsInfoMW)
	DashboardMultiplexer.Use(currentGuildDataMW)

	// Server manage pages
	RootMultiplexer.Handle(pat.New("/dashboard/:server/manage"), DashboardMultiplexer)
	RootMultiplexer.Handle(pat.New("/dashboard/:server/manage/*"), DashboardMultiplexer)

	DashboardMultiplexer.HandleFunc(pat.Get(""), EmbedHTML("manage.html"))
	DashboardMultiplexer.HandleFunc(pat.Get("/"), EmbedHTML("manage.html"))

	for _, route := range dashboardRoutes {
		route(DashboardMultiplexer)
	}
	return RootMultiplexer
}

func runRootMultiplexer() {
	mux := goji.NewMux()
	RootMultiplexer = mux

	mux.Use(baseTemplateDataMW)
	mux.Use(urlDataMW)

	mux.Handle(pat.Get("/static/*"), http.FileServer(http.FS(StaticFiles)))

	// Serve the login page
	mux.HandleFunc(pat.Get("/"), EmbedHTML("index.html"))
	mux.HandleFunc(pat.Get("/login"), handleLogin)
	mux.HandleFunc(pat.Get("/logout"), handleLogout)
	mux.HandleFunc(pat.Get("/confirm"), confirmLogin)

	// Data and service related pages
	mux.HandleFunc(pat.Get("/terms"), EmbedHTML("terms.html"))
	mux.HandleFunc(pat.Get("/privacy"), EmbedHTML("privacy.html"))
}

func runWebServer(multiplexer *goji.Mux) {
	logrus.Info("Webserver started on :80")
	http.ListenAndServe(":80", multiplexer)
}

// Called by plugins to add their routes
func RegisterDashboardRoutes(route func(*goji.Mux)) {
	dashboardRoutes = append(dashboardRoutes, route)
}

type FormResponse struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
}

func SendErrorToast(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(FormResponse{Success: false, Message: fmt.Sprintf("%s (ask support for help)", message)})
}
func SendSuccessToast(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(FormResponse{Success: true, Message: message})
}
