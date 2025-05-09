package web

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/RhykerWells/asbwig/frontend"
	"github.com/sirupsen/logrus"
	"goji.io/v3"
	"goji.io/v3/pat"
)

var (
	RootMultiplexer *goji.Mux
	HTMLTemplates fs.FS = frontend.HTMLTemplates
	StaticFiles fs.FS = frontend.StaticFiles
)

func Run() {
	initDiscordOauth()
	runRootMultiplexer()
	runWebServer(RootMultiplexer)
}

func embedHTML(filename string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func runRootMultiplexer() {
	mux := goji.NewMux()
	RootMultiplexer = mux

	mux.Handle(pat.Get("/static/*"), http.FileServer(http.FS(StaticFiles)))

	// Serve the login page
	mux.HandleFunc(pat.Get("/"), handleHomePage)
	mux.HandleFunc(pat.Get("/login"), handleLogin)
	mux.HandleFunc(pat.Get("/logout"), handleLogout)
	mux.HandleFunc(pat.Get("/confirm"), confirmLogin)
}

func runWebServer(multiplexer *goji.Mux) {
	logrus.Info("Webserver started on :8085")
	http.ListenAndServe(":8085", multiplexer)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	userData, _ := checkCookie(w, r)
	embedHTML("index.html", map[string]interface{}{"User": userData})(w, r)
}