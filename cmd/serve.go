package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"github.com/spf13/cobra"
)

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

		_, project, _, err := LoadCommon(cmd)
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

			// For now always generate a new doc to be able to reflect any change instantly.
			templates, project, extensions, err := LoadCommon(cmd)

			// For now the Generator is not necessarily thread save, so always create a new instance!
			var gen core.Generator = &generator.HTML{
				Markdown: generator.Markdown{
					DocTemplates: templates,
				},
			}
			atwhy, err := core.New(w, gen, project, extensions)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			if err := atwhy.Run(); err != nil {
				fmt.Println(err)
			}
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
