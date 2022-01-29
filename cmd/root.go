package cmd

import (
	"path/filepath"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// @WHY readme_usage
// Usage:
// Just run
// ```bash
// atwhy --help
// ```
// If nothing special is needed, just run the command without any arguments:
// ```bash
// atwhy
// ```
// It will use the default values and just work if a `templates` folder with some
// templates (e.g. `templates/README.tpl.md`) exists.

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atwhy",
	Short: "A documentation generator",
	Long: `atwhy can generate documentations based on templates.
It allows you to include documentation from anywhere in the project
and therefore provides a way to use "single source of truth" also for documentation.

Templates define how to combine the documentation annotations from all over the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateFolder, projectPath, extensions, err := LoadCommonArgs(cmd)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		var gen core.Generator
		generatorType, err := cmd.Flags().GetString("generator")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		switch generatorType {
		case "md":
			gen = generator.Markdown{}
		case "html":
			gen = &generator.HTML{
				Markdown: generator.Markdown{},
			}
		}

		atwhy, err := core.New(gen, projectPath, "", templateFolder, extensions)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		templates, err := atwhy.Load()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		for _, t := range templates {
			projectFS := afero.NewBasePathFs(afero.NewOsFs(), projectPath)

			err := projectFS.MkdirAll(t.Path, 0775)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			filename := filepath.Join(t.Path, t.Name+atwhy.Generator.Ext())
			file, err := projectFS.Create(filename)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			err = atwhy.Generate(t, file)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init is run by Go on startup. https://tutorialedge.net/golang/the-go-init-function/
func init() {
	// Global flags.
	rootCmd.PersistentFlags().StringP("templates-folder", "t", "templates", "path to a folder which contains the templates relative to the project directory")
	rootCmd.PersistentFlags().StringSliceP("ext", "e", nil, "comma separated list of allowed extensions\nallow all if not provided\nexample: .go,.js,.ts")
	rootCmd.PersistentFlags().StringP("project", "p", "", "the project folder")

	// Flags only for the root cmd.
	rootCmd.Flags().StringP("generator", "g", "md", "the generator to use\npossible values are: 'md', 'html'")
}
