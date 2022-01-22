package generator

import (
	mdTemplate "github.com/Tiffinger-Thiel-GmbH/atwhy/template"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"io"
	"strings"
)

type HTML struct {
	Markdown
}

func (h HTML) Ext() string {
	return ".html"
}

func (h HTML) Generate(markdownTemplate mdTemplate.Markdown, writer io.Writer) error {
	gm := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
		),
	)

	resMD := strings.Builder{}

	err := h.Markdown.Generate(markdownTemplate, &resMD)
	if err != nil {
		return err
	}

	return gm.Convert([]byte(resMD.String()), writer)
}
