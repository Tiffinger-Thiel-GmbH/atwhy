package tag

import "strconv"

type Basic struct {
	tagType Type
	value   string
}

func (b Basic) Type() Type {
	return b.tagType
}

func (b Basic) Markdown() string {
	return b.value
}

func (b Basic) IsParent() bool {
	return false
}

func FileLink(input Raw, _ []Tag) Tag {
	if input.Type != TypeFileLink {
		return nil
	}

	return Basic{
		tagType: input.Type,
		value:   "[" + input.Filename + ":" + strconv.Itoa(input.Line) + "](" + input.Filename + ")",
	}
}
