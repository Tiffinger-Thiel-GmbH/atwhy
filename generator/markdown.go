package generator

import (
	"io"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type Markdown struct {
	DocTemplates []DocTemplate
}

func createTagMap(tags []tag.Tag) map[string]tag.Tag {
	tagMap := make(map[string]tag.Tag)

	for _, t := range tags {
		tagMap[t.Placeholder()] = t
	}

	return tagMap
}

func (m *Markdown) Generate(tags []tag.Tag, writer io.Writer) error {
	tagMap := createTagMap(tags)

	for _, tpl := range m.DocTemplates {
		err := tpl.Execute(tagMap, writer)
		if err != nil {
			return err
		}

		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
