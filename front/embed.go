package front

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// https://www.liip.ch/en/blog/embed-sveltekit-into-a-go-binary
//
//go:generate npm i
//go:generate npm run build
//go:embed all:build
var files embed.FS

func StaticHandler(path string) http.Handler {
	fsys, err := fs.Sub(files, "build")
	if err != nil {
		log.Fatal(err)
	}
	filesystem := http.FS(fsys)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, path)
		// try if file exists at path, if not append .html (SvelteKit adapter-static specific)
		_, err := filesystem.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			path = fmt.Sprintf("%s.html", path)
		}
		r.URL.Path = path
		http.FileServer(filesystem).ServeHTTP(w, r)
	})
}
