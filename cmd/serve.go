package cmd

import (
	"fmt"
	"log"
	"net/http"

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

		_, project, _, err := LoadCommonArgs(cmd)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		// Serve the files.
		fs := http.FileServer(http.Dir(project))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				fs.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=UTF-8")

			/*	// For now always generate a new doc to be able to reflect any change instantly.
				templates, project, extensions, err := LoadCommonArgs(cmd)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				// For now the Generator is not necessarily thread save, so always create a new instance!
				var gen core.Generator = &generator.HTML{
					Markdown: generator.Markdown{},
				}
				atwhy, err := core.New(w, gen, project, extensions)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				if err := atwhy.Run(); err != nil {
					cmd.PrintErr(err)
					return
				}*/
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
