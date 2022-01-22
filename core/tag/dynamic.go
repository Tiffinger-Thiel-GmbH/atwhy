package tag

import (
	"path/filepath"
	"strconv"
	"strings"
)

type Dynamic struct {
	tagType     Type
	placeholder string
	factory     func(ctx Context) Tag
}

func (d Dynamic) Type() Type {
	return d.tagType
}

func (d Dynamic) String() string {
	return "Not evaluated yet. WithTemplate has not been used."
}

func (d Dynamic) Placeholder() string {
	return d.placeholder
}

func (d Dynamic) WithContext(ctx Context) Tag {
	return d.factory(ctx)
}

func Link(input Raw) (Tag, error) {
	if input.Type != TypeLink {
		return nil, nil
	}

	return Dynamic{
		tagType:     input.Type,
		placeholder: input.Placeholder,
		factory: func(ctx Context) Tag {
			// Add ../ for each part of the templatePath,
			// so that the relative path is still correct even if the destination path is in a sub folder.
			for _, part := range strings.Split(filepath.ToSlash(strings.TrimPrefix(ctx.TemplatePath(), ".")), "/") {
				if part != "" {
					input.Filename = "../" + input.Filename
				}
			}

			return Basic{
				tagType:     input.Type,
				placeholder: input.Placeholder,
				value:       "[" + input.Filename + ":" + strconv.Itoa(input.Line) + "](" + strings.TrimPrefix(input.Filename, ctx.TemplatePath()) + ")",
			}
		},
	}, nil
}
