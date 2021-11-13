package main

import (
	"fmt"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

/**
 * @FILELINE
 * @README
 * # My super Headline.
 * ##
 * __Body__ sdfag fd
 * dfgdfsg
 */

type Cleaner interface {
	Clean(string) (string, error)
}

type Processor struct {
	cleaners     []Cleaner
	tagFactories []tag.Factory
}

func (p Processor) clean(input string) string {
	for _, cleaner := range p.cleaners {
		cleaned, err := cleaner.Clean(input)
		if err != nil {
			// Just log out and continue with the other cleaners.
			fmt.Println(err)
			continue
		}
		input = cleaned
	}

	return input
}

func (p Processor) Process(tags []tag.Raw) ([]tag.Tag, error) {
	var processed []tag.Tag
	// Explicitly init the lastChildren to avoid nil.
	lastChildren := make([]tag.Tag, 0)
	var currentFile string

	for i := range tags {
		t := tags[i]
		t.Value = p.clean(t.Value)

		if currentFile != t.Filename {
			lastChildren = nil
		}
		currentFile = t.Filename

		for _, factory := range p.tagFactories {
			newTag := factory(t, lastChildren)
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
