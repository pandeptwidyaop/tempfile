// Package web embeds the static and template files for the web UI.
package web

import "embed"

// StaticFiles contains all static assets (CSS, JS).
//
//go:embed all:static
var StaticFiles embed.FS

// TemplateFiles contains all HTML templates.
//
//go:embed templates/*.html
var TemplateFiles embed.FS
