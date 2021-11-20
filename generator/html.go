package generator

import (
	"github.com/yuin/goldmark"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
	"io"
	"strings"
)

type HTML struct {
	Markdown
}

func (h HTML) Generate(tags []tag.Tag, writer io.Writer) error {
	gm := goldmark.New()

	var pages []string

	tagTypes := h.Markdown.TagsToExport

	for _, tagType := range tagTypes {
		resMD := strings.Builder{}
		h.Markdown.TagsToExport = []string{tagType}
		err := h.Markdown.Generate(tags, &resMD)
		if err != nil {
			return err
		}

		resHTML := strings.Builder{}
		err = gm.Convert([]byte(resMD.String()), &resHTML)
		if err != nil {
			return err
		}

		pages = append(pages, resHTML.String())
	}

	writer.Write([]byte(`<head>
    <meta charset="utf-8">
    <title>CrazyDoc</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
  </head>
  <body>`))

	writer.Write([]byte(`
		<nav class="navbar navbar-expand-lg navbar-light bg-light">
		<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarNavAltMarkup">
			<div class="navbar-nav">`))
	for _, tagType := range tagTypes {
		writer.Write([]byte(`
				<a class="nav-item nav-link active" onClick="document.querySelectorAll('.tagpage').forEach((p) => p.style.display = 'none'); document.querySelector('#` + tagType + `').style.display = 'block'" href="#">` + tagType + `</a>
			
`))
	}
	writer.Write([]byte(`</div>
	</div>
	</nav><div class="container">`))

	for i, page := range pages {

		page = `<div class="tagpage" id="` + string(tagTypes[i]) + `"/>` + page + `</div>`

		_, err := writer.Write([]byte(page))
		if err != nil {
			return err
		}
	}

	writer.Write([]byte(`</div></body>	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
	`))

	return nil
}
