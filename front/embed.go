package front

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// https://www.liip.ch/en/blog/embed-sveltekit-into-a-go-binary
//
//go:generate npm i
//go:generate npm run build
//go:embed all:build
var files embed.FS

// Get the subtree of the embedded files with `build` directory as a root.
func buildHTTPFS() http.FileSystem {
	build, err := fs.Sub(files, "build")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(build)
}

/*
*
https://dev.to/aryaprakasa/serving-single-page-application-in-a-single-binary-file-with-go-12ij#embed-the-static-files
*/
func HandleSPA() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buildPath := "build"
		f, err := files.Open(filepath.Join(buildPath, r.URL.Path))
		if os.IsNotExist(err) {
			index, err := files.ReadFile(filepath.Join(buildPath, "index.html"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(index)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		http.FileServer(buildHTTPFS()).ServeHTTP(w, r)
	})
}
