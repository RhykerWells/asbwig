package moderation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RhykerWells/asbwig/web"
	"goji.io/v3"
	"goji.io/v3/pat"
)

func initWeb() {
	web.RegisterDashboardRoutes(registerModerationRoutes)
}

func registerModerationRoutes(dashboard *goji.Mux) {
	moderationMux := goji.SubMux()

	moderationMux.Use(moderationConfigMW)

	dashboard.Handle(pat.New("/moderation"), moderationMux)
	dashboard.Handle(pat.New("/moderation/"), moderationMux)

	moderationMux.HandleFunc(pat.Get(""), web.EmbedHTML("moderation.html"))
	moderationMux.HandleFunc(pat.Get("/"), web.EmbedHTML("moderation.html"))
}

// moderationConfigMW provides middleware to parse the moderation config data to the template data
func moderationConfigMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		guildID := pat.Param(r, "server")

		config := GetConfig(guildID)

		tmplData, _ := ctx.Value(web.CtxKeyTmplData).(web.TmplContextData)
		tmplData["ModerationConfig"] = config

		ctx = context.WithValue(ctx, web.CtxKeyTmplData, tmplData)
		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}