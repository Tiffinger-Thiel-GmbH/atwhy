package main

import (
	"io"
	"sort"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

func MarkdownMapper(t tag.Tag) string {
	// TODO Regex !
	return strings.ReplaceAll(t.String(), "# ", "## ")
}

type MarkdownGenerator struct {
}

func (mG MarkdownGenerator) Generate(tags []tag.Tag, writer io.Writer) error {
	groupedTags := make(map[tag.Type][]tag.Tag)
	for _, t := range tags {
		groupedTags[t.Type()] = append(groupedTags[t.Type()], t)
	}

	for tagType, tagGroup := range groupedTags {

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

// check if really implements everything from Generator interface
var _ Generator = (*MarkdownGenerator)(nil)
