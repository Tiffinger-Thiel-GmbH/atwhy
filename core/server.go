package core

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func (a *AtWhy) ListenAndServe(host string) error {
	fs := http.FileServer(http.Dir(a.projectPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fast path: no need to generate everything if no html is requested.
		if !strings.HasSuffix(r.URL.Path, "/") && !strings.HasSuffix(r.URL.Path, a.Generator.Ext()) {
			fs.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path
		if strings.HasSuffix(path, "/") {
			path = path + "index.html"
		}

		// For now always load the templates and tags to be able to reflect any change instantly.
		// Maybe later a fs watcher could be used...
		templates, err := a.Load()
		if err != nil {
			// TODO use a logger
			fmt.Println(err)
			return
		}

		// Only generate the requested file.
		for _, t := range templates {
			if filepath.Join(t.Path, t.Name+a.Generator.Ext()) == path[1:] ||
				(t.Header.Server.Index && filepath.Join(t.Path, "index.html") == path[1:]) {
				// Found something
				w.Header().Set("Content-Type", "text/html; charset=UTF-8")

				err = a.Generate(t, w)
				if err != nil {
					// TODO use a logger
					fmt.Println(err)
					return
				}
				return
			}
		}

		// If still nothing is served use the file server.
		fs.ServeHTTP(w, r)
	})

	fmt.Printf("Starting server on %s\n", host)
	return http.ListenAndServe(host, nil)
}
