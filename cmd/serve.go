package cmd

import (
	"fmt"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// @WHY readme_usage_serve
// You can also serve the documentation on default host `localhost:4444` with:
// ```bash
// atwhy serve --ext .go
// ```
// For more information run `atwhy serve --help`

// serveCmd allows to serve the documentation using a html webserver.
var serveCmd = &cobra.Command{
	Use:   "serve [HOST (e.g. localhost:3333)]",
	Short: "Serves the documentation using a webserver.",
	Long: `Serves the documentation using a webserver.
It serves it on the given host 
(e.g. ":4444" to listen on all addresses, "localhost:4444" to listen only on localhost)`,
	Run: func(cmd *cobra.Command, args []string) {
		host := cmd.Flags().Arg(0)
		if host == "" {
			host = "localhost:4444"
		}

		templateFolder, projectPath, extensions, err := LoadCommonArgs(cmd)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		var gen core.Generator = &generator.HTML{
			Markdown: generator.Markdown{},
		}
		atwhy, err := core.New(gen, projectPath, templateFolder, extensions)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		// Serve the files.
		fs := http.FileServer(http.Dir(projectPath))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Fast path: no need to generate everything if no html is requested.
			if !strings.HasSuffix(r.URL.Path, atwhy.Generator.Ext()) {
				fs.ServeHTTP(w, r)
				return
			}

			// For now always load the templates and tags to be able to reflect any change instantly.
			// Maybe later a fs watcher could be used...
			templates, err := atwhy.Load()
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			// Only generate the requested file.
			for _, t := range templates {
				if filepath.Join(t.Path, t.Name+atwhy.Generator.Ext()) == r.URL.Path[1:] {
					// Found something
					w.Header().Set("Content-Type", "text/html; charset=UTF-8")

					err = atwhy.Generate(t, w)
					if err != nil {
						cmd.PrintErr(err)
						return
					}
					return
				}
			}

			// If still nothing is served use the file server.
			fs.ServeHTTP(w, r)
		})

		fmt.Printf("Starting server on %s\n", host)
		if err := http.ListenAndServe(host, nil); err != nil {
			log.Fatal(err)
		}
	},

	Args: cobra.MaximumNArgs(1),
}

// init is run by Go on startup. https://tutorialedge.net/golang/the-go-init-function/
func init() {
	rootCmd.AddCommand(serveCmd)
}
