package processor

import (
	"github.com/Tiffinger-Thiel-GmbH/atwhy/tag"
)

type Cleaner interface {
	Clean(string) (string, error)
}

type Processor struct {
	TagFactories []tag.Factory
}

func (p Processor) Process(tags []tag.Raw) ([]tag.Tag, error) {
	var processed []tag.Tag

	for i := range tags {
		t := tags[i]

		for _, factory := range p.TagFactories {
			newTag, err := factory(t)
			if err != nil {
				return nil, err
			}
			if newTag == nil {
				continue
			}

			processed = append(processed, newTag)
		}
	}

	return processed, nil
}
