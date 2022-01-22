package tag

import (
	"path/filepath"
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

func (b Basic) WithContext(_ Context) Tag {
	return b
}

func textFactory(input Raw, isMarkdown bool) Basic {
	// First remove windows line endings.
	input.Value = strings.ReplaceAll(input.Value, "\r\n", "\r")

	// Remove first line.
	splitted := strings.SplitN(input.Value, "\n", 2)
	var body string

	// If a body exists, use that. If not just leave the value empty.
	if isMarkdown && len(splitted) >= 2 {
		body = strings.ReplaceAll(splitted[1], "\n", "  \n")
	} else if !isMarkdown {
		body = splitted[1]
	}

	body = strings.TrimRight(body, " \n")
	return Basic{
		tagType:     input.Type,
		placeholder: input.Placeholder,
		value:       body,
	}
}

func Doc(input Raw) (Tag, error) {
	if input.Type != TypeDoc {
		return nil, nil
	}

	newTag := textFactory(input, true)

	// Inject hard newlines, as the linter in many languages strips away empty spaces.
	newTag.value = strings.ReplaceAll(newTag.value, "\r", HardNewLine)

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
