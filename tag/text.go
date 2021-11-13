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
	markdown := t.header + NewLine
	for _, child := range t.children {
		markdown += child.Markdown() + NewLine
	}

	return markdown + t.value
}

func (t Text) IsParent() bool {
	return true
}

func textFactory(input Raw, children []Tag) Tag {
	// Remove the first line. (It only contains the tag -> not interesting)
	parts := strings.SplitAfterN(input.Value, "\n", 3)[1:]

	var header string
	if len(parts) > 0 {
		header = parts[0]
	}

	if len(parts) > 1 {
		input.Value = parts[1]
	} else {
		input.Value = ""
	}

	return Text{
		Basic: Basic{
			tagType: input.Type,
			value:   input.Value,
		},
		header:   header,
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
