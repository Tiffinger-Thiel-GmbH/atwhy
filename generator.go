package main

import (
	"io"
	"strings"
)

func MarkdownMapper(t ProcessedTag) string {
	t.Value = strings.ReplaceAll(t.Value, "# ", "## ")
	result := strings.Split(t.Value, "\n")
	if t.Children != nil {
		for _, c := range t.Children {
			if c.Type == TagFileLine {
				result = append(result, "")
				copy(result[2:], result[1:])
				result[1] = c.Value
			}
		}
	}

	return strings.Join(result, "\n")
}

type Generate struct {
}

func (mG Generate) Generate(tags []ProcessedTag, writer io.Writer) error {
	for _, t := range tags {
		switch t.Type {
		case TagReadme:
			writer.Write([]byte(MarkdownMapper(t)))
			writer.Write([]byte("\n"))
		}
	}

	return nil
}

// check if really implements everything from Generator interface
var _ Generator = (*Generate)(nil)
