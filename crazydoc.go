package main

import (
	"io"

	"gitlab.com/tiffinger-thiel/crazydoc/loader"
)

// CrazyDoc combines all parts of the application.
//
// @CODE
// @FILELINK
// @WHY 10
// Architecture
// The idea behind CrazyDoc is to have several interfaces,
// each one for a small purpose.
// These interfaces are then implemented by concrete implementations.
//
// It is possible to replace or mock each part of the application at any time.
// You can use CrazyDoc as lib and just provide your own implementations.
//
// The interfaces are:
// * `Loader` loads files from a given path.
// * `TagFinder` reads the file and returns all lines which are part of a found
//   tag. It Does not process the raw lines.
// * `TagProcessor` processes the raw data from the `TagFinder` and generates
//   Tags out of them. It may also clean comment-chars and spaces and combine some tags.
// * `Generator` is responsible for postprocessing the tags and output the final file.
//   which it just writes to the passed `Writer`.
//
// So the workflow is:
// Loader -> TagFinder = tagList []tag.Raw
// tagList -> TagProcessor -> Generator -> Writer
type CrazyDoc struct {
	Loader    Loader
	Finder    loader.TagFinder
	Processor TagProcessor
	Generator Generator
	Writer    io.Writer
}

// @CODE_END

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
