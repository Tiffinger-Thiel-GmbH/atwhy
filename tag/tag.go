package tag

type Type string

// @DOC readme_tags
// You can use `\@DOC <placeholder_name>` and then use that placeholder in any template.
// There are also some special tags:
// * `\@DOC LINK <placeholder_name>` can be used to just add a link to the file where the tag is in.
// * `\@DOC CODE <placeholder_name>` can be used to reference any code.
//   It has to be closed by `\@DOC CODE_END`

var (
	TypeDoc     Type = "DOC"
	TypeLink    Type = "LINK"
	TypeCode    Type = "CODE"
	TypeCodeEnd Type = "CODE_END"
)

// Raw represents a not yet processed tag.
type Raw struct {
	Type        Type
	Placeholder string
	Filename    string
	Line        int
	Value       string
}

// Tag which was parsed from the code.
type Tag interface {
	Type() Type

	// String can be used to render the whole tag including the children.
	// The format may be chosen by the implementation, but it has to fit to the
	// other interface used implementations.
	String() string

	Placeholder() string
}

// Factory describes a function which can convert a Raw tag into a normal Tag.
// If the resulting Tag does not return nil on Tag.Children(), it means that the
// children have been consumed and the next factory call should get a new slice.
type Factory func(input Raw) (Tag, error)
