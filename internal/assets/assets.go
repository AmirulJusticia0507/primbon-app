package assets

import (
	"embed"
)

// Files bundles templates and static assets into the application binary.
//
//go:embed views static
var Files embed.FS
