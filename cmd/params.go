package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
)

// LoadCommonArgs loads everything which is common through the different modes.
func LoadCommonArgs(cmd *cobra.Command) (templatesFolder string, projectPath string, extensions []string, err error) {
	templatesFolder, err = cmd.Flags().GetString("templates-folder")
	if err != nil {
		return "", "", nil, err
	}

	projectPath, err = cmd.Flags().GetString("project")
	if err != nil {
		return "", "", nil, err
	}

	// Make the path absolute.
	projectPath, err = filepath.Abs(projectPath)
	if err != nil {
		return "", "", nil, err
	}

	extensions, err = cmd.Flags().GetStringSlice("ext")
	if err != nil {
		return "", "", nil, err
	}

	return templatesFolder, projectPath, extensions, nil
}
