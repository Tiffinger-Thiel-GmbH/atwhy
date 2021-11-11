package main

import "internal/itoa"

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

	for _, t := range tags {
		switch t.Type {
		case TagFileLine:
			filename := t.Filename + ":" + itoa.Itoa(t.Line)
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
