package statics

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var distFS embed.FS

// GetFS returns the http.FileSystem for the embedded static files.
// It serves from the "dist" root.
func GetFS() http.FileSystem {
	// Sub root to "dist"
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}
