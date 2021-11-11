package main

import "io"

type MarkdownGenerator struct {
}

func (mG MarkdownGenerator) Generate(tags []ProcessedTag) (io.Reader, error) {
	for _, t := range tags {
		switch t.Type {
		case TagReadme:
			// t.Value
		}
	}
	return nil, nil
}

// check if really implements everything from Generator interface
var _ Generator = (*MarkdownGenerator)(nil)
