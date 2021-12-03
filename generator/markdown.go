package generator

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

type Markdown struct {
	TagsToExport []string
}

func (m *Markdown) Generate(tags []tag.Tag, writer io.Writer) error {
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

		for _, tagToExport := range m.TagsToExport {
			if strings.Contains(tagToExport, string(tagType)) {
				foundTagsToExport = true
			}
		}

		if !foundTagsToExport {
			continue
		}

		_, err := writer.Write([]byte("# " + string(tagType) + "\n"))
		if err != nil {
			return err
		}

		for _, t := range tagGroup {
			_, err := writer.Write([]byte(MarkdownMapper(t) + "\n"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
