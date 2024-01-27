package ui

import (
	"embed"
)

// Content_static holds our static web server content.
//
//go:embed static
var Content_static embed.FS

// content_templates holds our templates for our web server.
//
//go:embed templates
var Content_templates embed.FS
