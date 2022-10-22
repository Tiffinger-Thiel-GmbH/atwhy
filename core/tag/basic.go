package tag

import (
	"path/filepath"
	"strconv"
	"strings"
)

const (
	HardNewLine = "  \n"
)

type Basic struct {
	tagType     Type
	value       string
	placeholder string
}

func (b Basic) Type() Type {
	return b.tagType
}

func (b Basic) String() string {
	return b.value
}

func (b Basic) Placeholder() string {
	return b.placeholder
}

func textFactory(input Raw, isMarkdown bool) Basic {
	// First remove windows line endings.
	input.Value = strings.ReplaceAll(input.Value, "\r\n", "\n")

	// Remove first line.
	splitted := strings.SplitN(input.Value, "\n", 2)
	var body string

	// If a body exists, use that. If not just leave the value empty.
	if isMarkdown && len(splitted) >= 2 {
		// Inject hard newlines, as the linter in many languages strips away empty spaces.
		body = strings.ReplaceAll(splitted[1], "\n", HardNewLine)
	} else if !isMarkdown && len(splitted) >= 2 {
		body = splitted[1]
	}

	body = strings.TrimRight(body, " \n")
	return Basic{
		tagType:     input.Type,
		placeholder: input.Placeholder,
		value:       body,
	}
}

func ProjectLink(input Raw) (Tag, error) {
	if input.Type != TypeLink {
		return nil, nil
	}

	escapedProjectFile := strings.ReplaceAll(input.Filename, `"`, `\"`)
	escapedProjectFile = strings.ReplaceAll(escapedProjectFile, `)`, `\)`)
	escapedTitle := strings.ReplaceAll(input.Filename, "[", `\[`)
	escapedTitle = strings.ReplaceAll(escapedTitle, "]", `\]`)

	return Basic{
		tagType:     input.Type,
		placeholder: input.Placeholder,
		// Insert the link-path as relative to be able to replace it in the final rendering based on the template path.
		value: "[" + escapedTitle + ":" + strconv.Itoa(input.Line) + `]({{ .Project "` + escapedProjectFile + `" }})`,
	}, nil
}

func Doc(input Raw) (Tag, error) {
	if input.Type != TypeDoc {
		return nil, nil
	}

	newTag := textFactory(input, true)

	return newTag, nil
}

func Code(input Raw) (Tag, error) {
	if input.Type != TypeCode {
		return nil, nil
	}

	newTag := textFactory(input, false)

	codeType := filepath.Ext(input.Filename)[1:]
	newTag.value = "```" + codeType + "\n" + newTag.value + "\n```\n"

	return newTag, nil
}
