package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"github.com/spf13/cobra"
)

// @WHY readme_usage
// Usage:
// Just run
// ```bash
// atwhy --help
// ```
// A common usage to for example generate this README.md is:
// ```bash
// atwhy --ext .go --templates README README.md
// ```

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atwhy [outputFile (e.g. README.md)]",
	Short: "A documentation generator",
	Long: `atwhy can generate documentations based on templates.
It allows you to include documentation from anywhere in the project
and therefore provides a way to use "single source of truth" also for documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		var writer io.Writer = os.Stdout

		out := cmd.Flags().Arg(0)

		if out != "" {
			file, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			defer file.Close()

			writer = file
		}

		templates, project, extensions, htmlTemplate, err := LoadCommon(cmd)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		var gen core.Generator
		switch filepath.Ext(out) {
		case ".md", "":
			gen = &generator.Markdown{
				DocTemplates: templates,
			}
		case ".html":
			gen = &generator.HTML{
				HTMLTemplate: htmlTemplate,
				Markdown: generator.Markdown{
					DocTemplates: templates,
				},
			}
		}

		atwhy, err := core.New(writer, gen, project, extensions)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if err := atwhy.Run(); err != nil {
			cmd.PrintErr(err)
			return
		}
	},

	Args: cobra.MaximumNArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init is run by Go on startup. https://tutorialedge.net/golang/the-go-init-function/
func init() {
	rootCmd.PersistentFlags().StringP("templates-folder", "t", "templates", "path to a folder which contains the templates relative to the project directory")
	rootCmd.PersistentFlags().StringSliceP("templates", "T", nil, "comma separated list of templates to generate\ngenerates all if omitted\nexample: README,WHY to generate README.tpl.md and WHY.tpl.md")
	rootCmd.PersistentFlags().StringSliceP("ext", "e", nil, "comma separated list of allowed extensions\nallow all if not provided\nexample: .go,.js,.ts")
	rootCmd.PersistentFlags().StringP("project", "p", "", "the project folder")
	rootCmd.PersistentFlags().String("html-template", "", "a .gohtml file which should be used instead of the embedded html template")
}
