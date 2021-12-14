package cmd

import (
	"path/filepath"

	"github.com/Tiffinger-Thiel-GmbH/AtWhy/generator"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// LoadCommon loads everything which is common through the different modes.
func LoadCommon(cmd *cobra.Command) (templates []generator.DocTemplate, project string, extentions []string, err error) {
	templatesFolder, err := cmd.Flags().GetString("templates-folder")
	if err != nil {
		return nil, "", nil, err
	}

	templatesAllowList, err := cmd.Flags().GetStringSlice("templates")
	if err != nil {
		return nil, "", nil, err
	}

	project, err = cmd.Flags().GetString("project")
	if err != nil {
		return nil, "", nil, err
	}

	templates, err = generator.LoadDocTemplates(afero.NewOsFs(), filepath.Join(project, templatesFolder), templatesAllowList)
	if err != nil {
		return nil, "", nil, err
	}

	extentions, err = cmd.Flags().GetStringSlice("ext")
	if err != nil {
		return nil, "", nil, err
	}

	return templates, project, extentions, nil
}
