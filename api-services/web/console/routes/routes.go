package web

import (
	"net/http"
	"os"
	"path/filepath"
)

func Mount(mux *http.ServeMux, distDir string) {
	fs := http.FileServer(http.Dir(distDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		if path != "/" {
			filePath := filepath.Join(distDir, path)
			info, err := os.Stat(filePath)
			if err == nil && !info.IsDir() {
				fs.ServeHTTP(w, r)
				return
			}
		}

		r.URL.Path = "/"
		fs.ServeHTTP(w, r)
	})
}
