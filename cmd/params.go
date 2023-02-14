package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/finder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultComments = []string{
	// @WHY CODE readme_comments_builtin
	"://,/*,*/", // All other files
	"sh:#",
	`py:#,""","""`,
	"cmd:REM",
	`cmd:::`,
	"vb:'",
	"html,xml:,<!--,-->",
	"lua:--,--[[,]]",
	"sql:--",
	"rb:#,=begin,=end",
	// @WHY CODE_END
}

var ErrInvalidCommentString = errors.New("comment configuration has to be like '{extList}:{lineComment}[,{blockStart},{blockEnd}]' (see --help)")
var ErrInvalidCommentStringMissingBlock = fmt.Errorf("either blockStart or blockEnd is missing - %w", ErrInvalidCommentString)

// LoadCommonArgs loads everything which is common through the different modes.
func LoadCommonArgs(cmd *cobra.Command) (templatesFolder string, projectPath string, extensions []string, commentConfig map[string]finder.CommentConfig, conf *viper.Viper, err error) {
	conf, err = initializeConfig(cmd)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	templatesFolder = conf.GetString("templates-folder")

	projectPath = conf.GetString("project")

	// Make the path absolute.
	projectPath, err = filepath.Abs(projectPath)
	if err != nil {
		return "", "", nil, nil, conf, err
	}

	extensions = conf.GetStringSlice("ext")
	comments := conf.GetStringSlice("comment")

	fmt.Println(comments)

	if len(comments) == 0 {
		comments = append(comments, "DEFAULT")
	}

	// Replace DEFAULT with the default comments.
	extendedComments := make([]string, 0, len(comments))
	for i, comment := range comments {
		if comment == "DEFAULT" {
			preComments := comments[:i]
			extendedComments = append(preComments, defaultComments...)
		} else {
			extendedComments = append(extendedComments, comment)
		}
	}

	comments = extendedComments

	commentConfig, err = generateCommentConfig(comments)
	if err != nil {
		return "", "", nil, nil, conf, err
	}

	return templatesFolder, projectPath, extensions, commentConfig, conf, nil
}

func generateCommentConfig(comments []string) (map[string]finder.CommentConfig, error) {
	commentConfig := make(map[string]finder.CommentConfig)

	for _, comment := range comments {
		split := strings.SplitN(comment, ":", 2)
		if len(split) != 2 {
			return nil, ErrInvalidCommentString
		}

		commaPlaceholder := string([]byte{1})

		cfgSplit := strings.Split(strings.ReplaceAll(split[1], `\,`, commaPlaceholder), ",")

		commentExtensions := strings.Split(split[0], ",")
		for _, ext := range commentExtensions {
			ext = "." + ext

			if _, ok := commentConfig[ext]; !ok {
				commentConfig[ext] = finder.CommentConfig{}
			}

			cfg := commentConfig[ext]

			if len(cfgSplit) >= 1 && cfgSplit[0] != "" {
				cfg.LineComment = append(commentConfig[ext].LineComment, cfgSplit[0])
			}

			if len(cfgSplit) >= 3 {
				if cfgSplit[1] == "" || cfgSplit[2] == "" {
					return nil, ErrInvalidCommentStringMissingBlock
				}

				cfg.BlockStart = append(commentConfig[ext].BlockStart, cfgSplit[1])
				cfg.BlockEnd = append(commentConfig[ext].BlockEnd, cfgSplit[2])
			} else if len(cfgSplit) == 2 {
				return nil, ErrInvalidCommentStringMissingBlock
			}

			// If no config was found, just don't add the cfg.
			if len(cfg.LineComment) > 0 || len(cfg.BlockStart) > 0 {
				// Replace back the commaPlaceholder
				cfg.LineComment = replaceInSlice(cfg.LineComment, commaPlaceholder, ",")
				cfg.BlockStart = replaceInSlice(cfg.BlockStart, commaPlaceholder, ",")
				cfg.BlockEnd = replaceInSlice(cfg.BlockEnd, commaPlaceholder, ",")
				commentConfig[ext] = cfg
			}
		}
	}

	return commentConfig, nil
}

func replaceInSlice(source []string, value string, replacement string) []string {
	var result []string
	for _, in := range source {
		result = append(result, strings.ReplaceAll(in, value, replacement))
	}
	return result
}
