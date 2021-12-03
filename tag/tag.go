package tag

type Type string

// @DOC readme.tags
// * Text
//   * `\@README`
//   * `\@WHY`
// * Modifiers (use them just before the text-tags)
//   * `\@FILELINK`
//   * `\@CODE` (use `\@CODE_END` after the code-block you want to include)

var (
	TypeWhy      Type = "WHY"
	TypeReadme   Type = "README"
	TypeFileLink Type = "FILELINK"
	TypeCode     Type = "CODE"
	TypeCodeEnd  Type = "CODE_END"
)

const (
	NewLine = "  \n"
)

// Raw represents a not yet processed tag.
type Raw struct {
	Type     Type
	Filename string
	Line     int
	Value    string
	Code     string
}

// Tag which was parsed from the code.
type Tag interface {
	Type() Type

	IsParent() bool

	// String can be used to render the whole tag including the children.
	// The format may be chosen by the implementation, but it has to fit to the
	// other interface used implementations.
	String() string

	// Position can be used to sort the tags.
	Position() int
}

// Factory describes a function which can convert a Raw tag into a normal Tag.
// If the resulting Tag does not return nil on Tag.Children(), it means that the
// children have been consumed and the next factory call should get a new slice.
type Factory func(input Raw, children []Tag) (Tag, error)
