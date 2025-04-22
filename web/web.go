package web

import (
	"io/fs"
	"net/http"

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

func embedHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := frontend.HTMLTemplates.ReadFile("templates/" + filename)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	}
}

func runRootMultiplexer() {
	mux := goji.NewMux()
	RootMultiplexer = mux

	// Serve the login page
	mux.HandleFunc(pat.Get("/"), embedHTML("index.html"))
	mux.HandleFunc(pat.Get("/login"), handleLogin)
	mux.HandleFunc(pat.Get("/logout"), handleLogout)
	mux.HandleFunc(pat.Get("/confirm"), confirmLogin)
}

func runWebServer(multiplexer *goji.Mux) {
	logrus.Info("Webserver started on :8085")
	http.ListenAndServe(":8085", multiplexer)
}