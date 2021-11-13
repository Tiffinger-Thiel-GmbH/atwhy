package tag

import "strconv"

type TagType string

var (
	TypeWhy      TagType = "WHY"
	TypeReadme   TagType = "README"
	TypeFileLink TagType = "FILELINK"
)

var (
	newLine = "  \n"
)

// Raw represents a not yet processed tag.
type Raw struct {
	Type     TagType
	Filename string
	Line     int
	Value    string
}

// Tag which was parsed from the code.
type Tag interface {
	Type() TagType

	// Children should return nil, if no children are possible.
	// If they are possible but there is none, it should return an empty slice.
	Children() []Tag

	IsParent() bool

	// Markdown can be used to render the whole tag including the children.
	Markdown() string
}

// Factory describes a function which can convert a Raw tag into a normal Tag.
// If the resulting Tag does not return nil on Tag.Children(), it means that the
// children have been consumed and the next factory call should get a new slice.
type Factory func(input Raw, children []Tag) Tag

type Basic struct {
	tagType TagType
	value   string
}

func (b Basic) Type() TagType {
	return b.Type()
}

func (b Basic) Children() []Tag {
	return nil
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
