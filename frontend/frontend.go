package frontend

import "embed"

//go:embed static/*
var StaticFiles embed.FS

//go:embed templates/*
var HTMLTemplates embed.FS

//go:embed pages/*
var HTMLPages embed.FS
