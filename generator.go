package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

func MarkdownMapper(t tag.Tag) string {
	// TODO Regex !

	fmt.Println("hallo")
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

	for tagType, tagGroup := range groupedTags {
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

	err = gm.Convert([]byte(res.String()), writer)
	return err
}
