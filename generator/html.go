package generator

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/tag"
	"github.com/aligator/checkpoint"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

//go:embed template
var EmbeddedTemplateFS embed.FS

type HTML struct {
	Markdown

	// HTMLTemplate the path to a .gohtml template to use instead of the embedded one.
	HTMLTemplate string

	sysFS    fs.FS
	template *template.Template
}

// loadTemplate either from the given HTMLTemplate path
// or if that is not set, from the embedded template.
func (h *HTML) loadTemplate() error {
	if h.template == nil {
		if h.HTMLTemplate != "" {
			var err error
			h.template, err = template.ParseFiles(h.HTMLTemplate)
			return checkpoint.From(err)
		}

		var err error
		h.sysFS, err = fs.Sub(EmbeddedTemplateFS, "template")
		if err != nil {
			return checkpoint.From(err)
		}

		h.template, err = template.ParseFS(h.sysFS, "index.gohtml")
		if err != nil {
			return checkpoint.From(err)
		}
	}

	return nil
}

func (h *HTML) Generate(tags []tag.Tag, writer io.Writer) error {
	err := h.loadTemplate()
	if err != nil {
		return checkpoint.From(err)
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

	for _, tpl := range h.Markdown.DocTemplates {
		resMD := strings.Builder{}
		h.Markdown.DocTemplates = []DocTemplate{tpl}
		err := h.Markdown.Generate(tags, &resMD)
		if err != nil {
			return checkpoint.From(err)
		}

		resHTML := strings.Builder{}
		err = gm.Convert([]byte(resMD.String()), &resHTML)
		if err != nil {
			return checkpoint.From(err)
		}

		data.Pages = append(data.Pages, Page{
			ID:    tpl.ID,
			Title: tpl.Header.Meta.Title,
			Body:  template.HTML(resHTML.String()),
		})
	}

	tplName := filepath.Base(h.HTMLTemplate)
	if tplName == "" || tplName == "." {
		// index.gohtml is the name of the embedded template.
		tplName = "index.gohtml"
	}

	return checkpoint.From(h.template.ExecuteTemplate(writer, tplName, data))
}
