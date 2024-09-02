package view

import "embed"

var (
	//go:embed assets
	Assets embed.FS
	//go:embed templates
	Templates embed.FS
)
