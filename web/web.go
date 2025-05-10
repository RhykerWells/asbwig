package web

import (
	"io/fs"
	"net/http"
	"text/template"
	"time"

	"github.com/RhykerWells/asbwig/common"
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

	// Data and service related pages
	mux.HandleFunc(pat.Get("/terms"), handleTerms)
	mux.HandleFunc(pat.Get("/privacy-policy"), handlePrivacy)
}

func runWebServer(multiplexer *goji.Mux) {
	dashboard(multiplexer)

	logrus.Info("Webserver started on :8085")
	http.ListenAndServe(":8085", multiplexer)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	userData, _ := checkCookie(w, r)
	embedHTML("index.html", map[string]interface{}{"User": userData, "Year": time.Now().UTC().Year()})(w, r)
}

func handleTerms(w http.ResponseWriter, r *http.Request) {
	embedHTML("terms.html", map[string]interface{}{})(w,r)
}

func handlePrivacy(w http.ResponseWriter, r *http.Request) {
	embedHTML("privacy.html", map[string]interface{}{})(w,r)
}


func dashboard(mux *goji.Mux) {
	mux.HandleFunc(pat.Get("/dashboard"), handleDashboard)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	userData, _ := checkCookie(w, r)
    
    // Check that userData["id"] exists and is a string
    userID, _ := userData["id"].(string)
	// Retrieve the guilds managed by the user
	guilds := getUserManagedGuilds(common.Session, userID)
	// Create a map to store guild data (ID and Name)
	guildList := make([]map[string]interface{}, 0)
	for guildID, guildName := range guilds {
		avatarURL := "./static/img/icons/cross.png"
		if guild, err := common.Session.Guild(guildID); err == nil {
			if url := guild.IconURL("1024"); url != "" {
				avatarURL = url
			}
		}
		guildList = append(guildList, map[string]interface{}{
			"ID":   guildID,
			"Avatar": avatarURL,
			"Name": guildName,
		})
	}

	// Marshal the guild data into JSON and write to the response
	responseData := map[string]interface{}{
		"User": userData,
		"Guilds": guildList,
	}

	embedHTML("dashboard.html", responseData)(w,r)
}