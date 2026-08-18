package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Mount(mux *http.ServeMux, distDir string) {
	fs := http.FileServer(http.Dir(distDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")

		if path != "" {
			filePath := filepath.Join(distDir, path)

			info, err := os.Stat(filePath)
			if err == nil && !info.IsDir() {
				r.URL.Path = "/" + path
				fs.ServeHTTP(w, r)
				return
			}
		}

		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	})
}
