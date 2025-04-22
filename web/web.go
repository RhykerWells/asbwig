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
)

func Run() {
	initDiscordOauth()
	runRootMultiplexer()
	runWebServer(RootMultiplexer)
}

func embedHTML(filename string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileData, err := frontend.HTMLTemplates.ReadFile("templates/" + filename)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.New(filename).Parse(string(fileData))
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}

func runRootMultiplexer() {
	mux := goji.NewMux()
	RootMultiplexer = mux

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