package web

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"
	"sync"
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

	coreHTMLTemplates   fs.FS = frontend.HTMLTemplates
	coreHTMLPages       fs.FS = frontend.HTMLPages
	additionalHTMLPages []fs.FS
	tmplOnce            sync.Once
	tmpl                *template.Template
	StaticFiles         fs.FS = frontend.StaticFiles

	URL        string = "https://" + common.ConfigASBWIGHost
	TermsURL   string = common.ConfigTermsURLOverride
	PrivacyURL string = common.ConfigPrivacyURLOverride

	dashboardRoutes []func(*goji.Mux)
)

// Run starts the necessary processes to begin serving the web server
// it handles starting the authentication and the multiplexer
func Run() {
	initDiscordOauth()
	multiplexer := setupWebRoutes()
	runWebServer(multiplexer)
}

func AddHTMLFilesystem(fs fs.FS) {
	additionalHTMLPages = append(additionalHTMLPages, fs)
}

func loadTemplates() (*template.Template, error) {
	var tmplError error

	tmplOnce.Do(func() {
		// Set the initial template with the template functions
		t := template.New("").Funcs(templateFunctions)

		// Parse core templates
		_, err := t.ParseFS(coreHTMLTemplates, "templates/*.html")
		if err != nil {
			tmplError = err
			return
		}

		// Parse core pages
		_, err = t.ParseFS(coreHTMLPages, "pages/*.html")
		if err != nil {
			tmplError = err
			return
		}

		for _, fsys := range additionalHTMLPages {
			files, err := fs.Glob(fsys, "*/*.html")
			if err != nil {
				tmplError = err
				return
			}

			for _, file := range files {
				data, err := fs.ReadFile(fsys, file)
				if err != nil {
					tmplError = err
					return
				}
				if _, err := t.New(filepath.Base(file)).Parse(string(data)); err != nil {
					tmplError = err
					return
				}
			}
		}
		tmpl = t
	})

	return tmpl, tmplError
}

// RenderPage serves the given html page
// It attaches the template function array and the template data
func RenderPage(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := loadTemplates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve the template data from the context
		tmplData, _ := r.Context().Value(CtxKeyTmplData).(TmplContextData)

		// Attempt to render the HTML templates and pages
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.ExecuteTemplate(w, filename, tmplData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// setupWebRoutes calls the initial root multiplexer and initialises the base dashboard web routes,
// and attaches the required middlewares for validation and template data
func setupWebRoutes() *goji.Mux {
	// start the base routes, such as logins, static files and such
	runRootMultiplexer()

	// Handle dashboard page
	RootMultiplexer.Handle(pat.Get("/dashboard"), userAndManagedGuildsInfoMW(RenderPage("dashboard.html")))
	RootMultiplexer.Handle(pat.Get("/dashboard/"), userAndManagedGuildsInfoMW(RenderPage("dashboard.html")))

	// Create a sub-mux for dashboard-related routes
	DashboardMultiplexer = goji.SubMux()

	// Middlewares
	DashboardMultiplexer.Use(validateGuild)
	DashboardMultiplexer.Use(userAndManagedGuildsInfoMW)
	DashboardMultiplexer.Use(currentGuildDataMW)

	// Server manage pages
	RootMultiplexer.Handle(pat.New("/dashboard/:server/manage"), DashboardMultiplexer)
	RootMultiplexer.Handle(pat.New("/dashboard/:server/manage/*"), DashboardMultiplexer)

	DashboardMultiplexer.HandleFunc(pat.Get(""), RenderPage("manage.html"))
	DashboardMultiplexer.HandleFunc(pat.Get("/"), RenderPage("manage.html"))

	// Attach additional dashboard sub routes to the multiplexer
	for _, route := range dashboardRoutes {
		route(DashboardMultiplexer)
	}
	return RootMultiplexer
}

// runRootMultiplexer creates the initial multiplexer and attaches the base middlewares for template data
func runRootMultiplexer() {
	// Create the mux
	mux := goji.NewMux()
	RootMultiplexer = mux

	// Middlewares
	mux.Use(baseTemplateDataMW)
	mux.Use(urlDataMW)

	// Static files
	mux.Handle(pat.Get("/static/*"), http.FileServer(http.FS(StaticFiles)))

	// Serve the login page
	mux.HandleFunc(pat.Get("/"), RenderPage("index.html"))
	mux.HandleFunc(pat.Get("/login"), handleLogin)
	mux.HandleFunc(pat.Get("/logout"), handleLogout)
	mux.HandleFunc(pat.Get("/confirm"), confirmLogin)

	// Data and service related pages
	mux.HandleFunc(pat.Get("/terms"), RenderPage("terms.html"))
	mux.HandleFunc(pat.Get("/privacy"), RenderPage("privacy.html"))
}

// runWebServer serves the multiplexer on the default http port
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

// SendSuccessToast is used to send a JSON response to the client to send the success toasts
func SendSuccessToast(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(FormResponse{Success: true, Message: message})
}

// SendErrorToast is used to send a JSON response to the client to send the error toasts
func SendErrorToast(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(FormResponse{Success: false, Message: message})
}
