package processor

import (
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type Cleaner interface {
	Clean(string) (string, error)
}

type Processor struct {
	TagFactories []tag.Factory
}

func (p Processor) Process(tags []tag.Raw) ([]tag.Tag, error) {
	var processed []tag.Tag
	// Explicitly init the lastChildren to avoid nil.
	lastChildren := make([]tag.Tag, 0)
	var currentFile string

	for i := range tags {
		t := tags[i]

		if currentFile != t.Filename {
			lastChildren = nil
		}
		currentFile = t.Filename

		for _, factory := range p.TagFactories {
			newTag, err := factory(t, lastChildren)
			if err != nil {
				return nil, err
			}
			if newTag == nil {
				continue
			}

			// If the new tag is a parent, it consumed all children.
			// This means all lastChildren are used now and the lastChildren can be reset to an empty slice.
			if newTag.IsParent() {
				lastChildren = []tag.Tag{}
				processed = append(processed, newTag)
			} else {
				// If the new tag is a child, it has to be cached for later.
				lastChildren = append(lastChildren, newTag)
			}
		}
	}

	return processed, nil
}
