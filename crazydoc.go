package main

import (
	"io"

	"gitlab.com/tiffinger-thiel/crazydoc/loader"
)

// CrazyDoc combines all parts of the application.
// @DOC LINK crazydoc.struct.link
// @DOC CODE crazydoc.struct.code
type CrazyDoc struct {
	Loader    Loader
	Finder    loader.TagFinder
	Processor TagProcessor
	Generator Generator
	Writer    io.Writer
}

// @DOC CODE_END

func (c CrazyDoc) Run(path string) error {
	tags, err := c.Loader.Load(path, c.Finder)
	if err != nil {
		return err
	}

	processed, err := c.Processor.Process(tags)
	if err != nil {
		return err
	}

	return c.Generator.Generate(processed, c.Writer)
}
