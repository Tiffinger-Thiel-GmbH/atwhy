package cmd

import (
	"path/filepath"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// @WHY readme_usage0
// __Generate__
// If nothing special is needed, just run the command without any arguments:
// ```bash
// atwhy
// ```
// It will use the default values and just work if a `templates` folder with some
// templates (e.g. `templates/README.tpl.md`) exists.
// For more information run `atwhy --help`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atwhy",
	Short: "A documentation generator",
	Long: `atwhy can generate documentations based on templates.
It allows you to include documentation from anywhere in the project
and therefore provides a way to use "single source of truth" also for documentation.

Templates define how to combine the documentation annotations from all over the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateFolder, projectPath, extensions, commentConfig, err := LoadCommonArgs(cmd)
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

		atwhy, err := core.New(gen, projectPath, "/", templateFolder, extensions, commentConfig)
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

	// @WHY readme_comments
	// Each `--comment` is a string with the following format:
	// `--comment={extList}:{lineComment},{blockStart},{blockEnd}`
	// Where:
	// * `extList` is a comma-separated list of file extensions,
	// * `lineComment` is the comment prefix for line comments, (e.g. `//` or `#`),
	// * `blockStart` is the comment prefix for block comments start, (e.g. `/*` or `<!--`),
	// * `blockEnd` is the comment prefix for block comments end, (e.g. `*/` or `-->`).
	// * escape ',' by '\,' if needed.
	// * A catch-all (e.g. "://,/*,*/") catches all not otherwise configured extensions.
	// __`blockStart` and `blockEnd` are optional.__
	//
	// Examples:
	// * set only the lineComment for sh files
	//   `"sh:#"`
	// * set only the blockComment for html,xml
	//   `"html,xml:,<!--,-->"`
	// * set c-style for all files (not caught by another rule before)
	//   `"://,/*,*/"`
	//
	// If `--comment` is passed at least one time, all built-in rules are disabled.
	// Use `--comment=DEFAULT` if you still want to use the built-in rules.

	rootCmd.PersistentFlags().StringArray("comment", nil, `Set the comments for a specific file ending.
Syntax: 
{extList}:{lineComment},{blockStart},{blockEnd}

extList consists of the comma-separated list of file extensions which match the specified comment configuration.
blockStart and blockEnd may be omitted if only lineComments exist.
If no lineComment exists, just leave it blank.
If a file extension is in more than one rules, all are added to that extension.
A catch-all (e.g. "://,/*,*/") catches all not otherwise configured extensions.

Escape ',' by '\,' if needed.

Examples:
* set only the lineComment for sh files
  "sh:#"
* set only the blockComment for html,xml
  "html,xml:,<!--,-->"
* set c-style for all files (not caught by another rule before)
  "://,/*,*/" 

If "--comment" is passed at least one time, all built-in rules are disabled.
Use "--comment=DEFAULT" if you still want to use the built-in rules.
`)
	// Flags only for the root cmd.
	rootCmd.Flags().StringP("generator", "g", "md", "the generator to use\npossible values are: 'md', 'html'")
}
