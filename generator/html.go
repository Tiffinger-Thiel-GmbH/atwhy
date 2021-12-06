package generator

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

//go:embed template
var TemplateFS embed.FS

type HTML struct {
	Markdown

	template *template.Template
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

func (h *HTML) Generate(tags []tag.Tag, writer io.Writer) error {
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
		ID   string
		Name string
		Body template.HTML
	}
	type Data struct {
		Pages []Page
	}

	var data Data

	templates := h.Markdown.DocTemplates

	for _, tpl := range templates {
		resMD := strings.Builder{}
		h.Markdown.DocTemplates = []DocTemplate{tpl}
		err := h.Markdown.Generate(tags, &resMD)
		if err != nil {
			return err
		}

		resHTML := strings.Builder{}
		err = gm.Convert([]byte(resMD.String()), &resHTML)
		if err != nil {
			return err
		}

		data.Pages = append(data.Pages, Page{
			ID:   tpl.ID,
			Name: tpl.Name,
			Body: template.HTML(resHTML.String()),
		})
	}

	return h.template.ExecuteTemplate(writer, "index.gohtml", data)
}
