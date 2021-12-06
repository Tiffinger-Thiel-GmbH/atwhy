package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/tiffinger-thiel/crazydoc/core"
	"gitlab.com/tiffinger-thiel/crazydoc/generator"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "With serve you can start a webserver which serves the documentation",
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

			// For now the Generator is not necessarily thread save!
			var gen core.Generator = &generator.HTML{
				Markdown: generator.Markdown{
					DocTemplates: templates,
				},
			}
			crazyDoc, err := core.New(w, gen, project, extensions)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			if err := crazyDoc.Run(); err != nil {
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

func init() {
	rootCmd.AddCommand(serveCmd)
}
