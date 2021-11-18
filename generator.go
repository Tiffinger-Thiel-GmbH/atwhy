package main

import (
	"io"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

func MarkdownMapper(t tag.Tag) string {
	// TODO Regex !

	return strings.ReplaceAll(t.String(), "# ", "## ")
}

type MarkdownGenerator struct {
	TagsToExport []string
}

func (mg MarkdownGenerator) Generate(tags []tag.Tag, writer io.Writer) error {
	groupedTags := make(map[tag.Type][]tag.Tag)
	for _, t := range tags {
		groupedTags[t.Type()] = append(groupedTags[t.Type()], t)
	}

	var sorted []tag.Type
	for tagType := range groupedTags {
		sorted = append(sorted, tagType)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	for _, tagType := range sorted {
		tagGroup := groupedTags[tagType]
		foundTagsToExport := false

		for _, tagToExport := range mg.TagsToExport {
			if strings.Contains(tagToExport, string(tagType)) {
				foundTagsToExport = true
			}
		}

		if !foundTagsToExport {
			continue
		}

		sort.Slice(tagGroup, func(i, j int) bool {
			return tagGroup[i].Position() < tagGroup[j].Position()
		})

		writer.Write([]byte("# " + string(tagType) + "\n"))
		for _, tag := range tagGroup {
			writer.Write([]byte(MarkdownMapper(tag) + "\n\n"))
		}
	}

	return nil
}

type HTMLGenerator struct {
	MarkdownGenerator
}

func (hg HTMLGenerator) Generate(tags []tag.Tag, writer io.Writer) error {
	res := strings.Builder{}
	err := hg.MarkdownGenerator.Generate(tags, &res)
	if err != nil {
		return err
	}

	gm := goldmark.New()

	writer.Write([]byte(`<head>
    <meta charset="utf-8">
    <title>CrazyDoc</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">

    <style>
      
    </style>
  </head>
  <body>`))

	err = gm.Convert([]byte(res.String()), writer)

	writer.Write([]byte(`</body>`))

	return err
}
