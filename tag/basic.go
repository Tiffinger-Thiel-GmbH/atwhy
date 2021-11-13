package tag

import "strconv"

type Basic struct {
	tagType Type
	value   string
}

func (b Basic) Type() Type {
	return b.tagType
}

func (b Basic) String() string {
	return b.value
}

func (b Basic) IsParent() bool {
	return false
}

func (b Basic) Position() int {
	return 0
}

func FileLink(input Raw, _ []Tag) (Tag, error) {
	if input.Type != TypeFileLink {
		return nil, nil
	}

	return Basic{
		tagType: input.Type,
		value:   "[" + input.Filename + ":" + strconv.Itoa(input.Line) + "](" + input.Filename + ")",
	}, nil
}
