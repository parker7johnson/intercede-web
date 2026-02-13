package web

import "embed"

// StaticFiles embeds the static directory for production builds
//
//go:embed static/*
var StaticFiles embed.FS
