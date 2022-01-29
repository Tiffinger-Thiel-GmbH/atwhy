package cmd

import (
	"log"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"github.com/spf13/cobra"
)

// @WHY readme_usage_serve
// You can also serve the documentation on default host `localhost:4444` with:
// ```bash
// atwhy serve
// ```
// For more information run `atwhy serve --help`

// serveCmd allows to serve the documentation using a html webserver.
var serveCmd = &cobra.Command{
	Use:   "serve [HOST (e.g. localhost:4444)]",
	Short: "Serves the documentation using a webserver.",
	Long: `Serves the documentation using a webserver.
It serves it on the given host. 
(e.g. ":4444" to listen on all addresses, "localhost:4444" to listen only on localhost)
Default is: "localhost:4444"`,
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
		atwhy, err := core.New(gen, projectPath, "/project/", templateFolder, extensions)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if err := atwhy.ListenAndServe(host); err != nil {
			log.Fatal(err)
		}
	},

	Args: cobra.MaximumNArgs(1),
}

// init is run by Go on startup. https://tutorialedge.net/golang/the-go-init-function/
func init() {
	rootCmd.AddCommand(serveCmd)
}
