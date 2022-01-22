package generator

import (
	"embed"
	mdTemplate "github.com/Tiffinger-Thiel-GmbH/atwhy/template"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

//go:embed template
var TemplateFS embed.FS

type HTML struct {
	Markdown

	template *template.Template
}

func (h HTML) Ext() string {
	return ".html"
}

func (h *HTML) loadTemplate() error {
	if h.template == nil {
		subFS, err := fs.Sub(TemplateFS, "template")
		if err != nil {
			return err
		}

		h.template, err = template.ParseFS(subFS, "index.gohtml")
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HTML) Generate(markdownTemplate mdTemplate.MarkdownTemplate, writer io.Writer) error {
	err := h.loadTemplate()
	if err != nil {
		return err
	}

	gm := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
		),
	)

	type Page struct {
		ID    string
		Title string
		Body  template.HTML
	}
	type Data struct {
		Pages []Page
	}

	var data Data

	resMD := strings.Builder{}

	err = h.Markdown.Generate(markdownTemplate, &resMD)
	if err != nil {
		return err
	}

	resHTML := strings.Builder{}
	err = gm.Convert([]byte(resMD.String()), &resHTML)
	if err != nil {
		return err
	}

	data.Pages = append(data.Pages, Page{
		ID:    markdownTemplate.ID,
		Title: markdownTemplate.Header.Meta.Title,
		Body:  template.HTML(resHTML.String()),
	})

	return h.template.ExecuteTemplate(writer, "index.gohtml", data)
}
