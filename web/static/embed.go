package static

import "embed"

//go:embed */**
var StaticDir embed.FS

