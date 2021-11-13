package tag

import (
	"regexp"
	"strconv"
	"strings"
)

type Text struct {
	Basic
	header   string
	children []Tag
	position int
}

func (t Text) String() string {
	var markdown string

	if strings.Trim(t.header, " \n") != "" {
		markdown = "# " + t.header
	}

	for _, child := range t.children {
		markdown += child.String() + NewLine
	}

	return markdown + t.value
}

func (t Text) IsParent() bool {
	return true
}

func (t Text) Position() int {
	return t.position
}

func textFactory(input Raw, children []Tag) (Tag, error) {
	// Split into 1 Tag-line, 1 header-line and the rest.
	parts := strings.SplitAfterN(input.Value, "\n", 3)

	// Read the position argument.
	reg, err := regexp.Compile("@" + string(input.Type) + " (\\d+)")
	if err != nil {
		return nil, err
	}

	found := reg.FindStringSubmatch(parts[0])
	var position int
	if len(found) >= 2 {
		position, err = strconv.Atoi(found[1])
		if err != nil {
			return nil, err
		}
	}

	// Remove the first line. (It only contains the tag -> not interesting)
	parts = parts[1:]

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
		position: position,
	}, nil
}

func Why(input Raw, children []Tag) (Tag, error) {
	if input.Type != TypeWhy {
		return nil, nil
	}

	return textFactory(input, children)
}

func Readme(input Raw, children []Tag) (Tag, error) {
	if input.Type != TypeReadme {
		return nil, nil
	}

	return textFactory(input, children)
}
