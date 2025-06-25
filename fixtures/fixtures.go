package fixtures

import (
	"embed"
)

//go:embed data/*
var FS embed.FS
