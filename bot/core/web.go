package core

import (
	"context"
	"net/http"

	"github.com/RhykerWells/asbwig/web"
	"goji.io/v3"
	"goji.io/v3/pat"
)

func initWeb() {
	web.RegisterDashboardRoutes(registerCoreRoute)
}

func registerCoreRoute(dashboard *goji.Mux) {
	coreMux := goji.SubMux()

	coreMux.Use(coreConfigMW)

	dashboard.Handle(pat.New("/core"), coreMux)
	dashboard.Handle(pat.New("/core/"), coreMux)

	coreMux.HandleFunc(pat.Get(""), web.EmbedHTML("core.html"))
	coreMux.HandleFunc(pat.Get("/"), web.EmbedHTML("core.html"))

	coreMux.HandleFunc(pat.Post(""), saveConfigHandler)
	coreMux.HandleFunc(pat.Post("/"), saveConfigHandler)
}

func saveConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	guildID := pat.Param(r, "server")
	config := GetConfig(guildID)

	formType := r.FormValue("form_type")
	switch formType {
	case "Core":
		if r.FormValue("GuildPrefix") == "" {
			web.SendErrorToast(w, "Prefix cannot be empty.")
			return
		}
		config.GuildPrefix = r.FormValue("GuildPrefix")
	}

	err := SaveConfig(config)
	if err != nil {
		web.SendErrorToast(w, err.Error())
		return
	}

	web.SendSuccessToast(w, "Successfully saved")
}

// coreConfigMW provides middleware to parse the core config data to the template data
func coreConfigMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		guildID := pat.Param(r, "server")

		config := GetConfig(guildID)

		tmplData, _ := ctx.Value(web.CtxKeyTmplData).(web.TmplContextData)
		tmplData["CoreConfig"] = config

		ctx = context.WithValue(ctx, web.CtxKeyTmplData, tmplData)
		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}