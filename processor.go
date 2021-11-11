package main

import "strconv"

/**
 * @FILELINE
 * @README
 * # My super Headline.
 * ##
 * __Body__ sdfag fd
 * dfgdfsg
 */

type Processor struct{}

func (p Processor) Process(tags []Tag) ([]ProcessedTag, error) {
	var processed []ProcessedTag
	var lastChildren []ProcessedTag
	var currentFile string

	for _, t := range tags {
		if currentFile != t.Filename {
			lastChildren = nil
		}
		currentFile = t.Filename

		switch t.Type {
		case TagFileLine:
			filename := t.Filename + ":" + strconv.Itoa(t.Line)
			lastChildren = append(lastChildren, ProcessedTag{
				Type:  t.Type,
				Value: "[" + filename + "](" + filename + ")",
			})
		default:
			processed = append(processed, ProcessedTag{
				Type:     t.Type,
				Value:    t.Value,
				Children: lastChildren,
			})
		}
	}

	return processed, nil
}
