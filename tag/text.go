package tag

import "strings"

type Text struct {
	Basic
	header   string
	children []Tag
}

func (t Text) Children() []Tag {
	return t.children
}

func (t Text) Markdown() string {
	markdown := t.header + newLine
	for _, child := range t.children {
		markdown += child.Markdown() + newLine
	}

	return markdown + t.value
}

func (t Text) IsParent() bool {
	return true
}

func textFactory(input Raw, children []Tag) Tag {
	// TODO: implement with logic that strips out the header

	// Remove the first line.
	parts := strings.SplitAfterN(input.Value, "\n", 2)[1:]
	if len(parts) > 0 {
		input.Value = parts[0]
	} else {
		input.Value = ""
	}

	return Text{
		Basic: Basic{
			tagType: input.Type,
			value:   input.Value,
		},
		header:   "",
		children: children,
	}
}

func Why(input Raw, children []Tag) Tag {
	if input.Type != TypeWhy {
		return nil
	}

	return textFactory(input, children)
}

func Readme(input Raw, children []Tag) Tag {
	if input.Type != TypeReadme {
		return nil
	}

	return textFactory(input, children)
}
