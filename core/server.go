package core

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	mdTemplate "github.com/Tiffinger-Thiel-GmbH/atwhy/template"
)

// Page currently just consists of a markdown template.
// So all information included can be used in the html template.
type Page = mdTemplate.Markdown

//go:embed html
var TemplateFS embed.FS

const templateFile = "page.gohtml"

func (a *AtWhy) initPageTemplate() error {
	subFS, err := fs.Sub(TemplateFS, "html")
	if err != nil {
		return err
	}

	a.pageTemplate, err = template.ParseFS(subFS, templateFile)
	return err
}

func (a *AtWhy) buildPage(writer io.Writer, pageID string, pages []Page) error {
	data := struct {
		ID    string
		Title string
		Body  template.HTML

		Pages []Page
	}{
		ID:    pageID,
		Pages: pages,
	}

	for _, page := range pages {
		if page.ID == pageID {
			buf := bytes.NewBufferString("")
			err := a.Generate(page, buf)
			if err != nil {
				return err
			}
			data.Body = template.HTML(buf.String())
			data.Title = page.Header.Meta.Title
		}
	}

	return a.pageTemplate.ExecuteTemplate(writer, templateFile, data)

}

func (a *AtWhy) ListenAndServe(host string) error {
	fileServer := http.StripPrefix(a.projectPathPrefix, http.FileServer(http.Dir(a.projectPath)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fast path: no need to generate everything if no html is requested.
		if strings.HasPrefix(r.URL.Path, a.projectPathPrefix) {
			fileServer.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path
		if strings.HasSuffix(path, "/") {
			path = path + "index.html"
		}

		// For now always load the templates and tags to be able to reflect any change instantly.
		// Maybe later a fs watcher could be used...
		templates, err := a.Load()
		if err != nil {
			// TODO use a logger
			fmt.Println(err)
			return
		}

		// Only generate the requested file.
		for _, t := range templates {
			if filepath.Join(t.Path, t.Name+a.Generator.Ext()) == path[1:] ||
				(t.Header.Server.Index && filepath.Join(t.Path, "index.html") == path[1:]) {
				// Found something
				w.Header().Set("Content-Type", "text/html; charset=UTF-8")

				err = a.buildPage(w, t.ID, templates)
				if err != nil {
					// TODO use a logger
					fmt.Println(err)
					return
				}
				return
			}
		}

		w.WriteHeader(404)
	})

	fmt.Printf("Starting server on %s\n", host)
	return http.ListenAndServe(host, nil)
}
