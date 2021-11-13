package main

import (
	"io"
)

type CrazyDoc struct {
	Finder    TagFinder
	Loader    Loader
	Processor TagProcessor
	Generator Generator
	Writer    io.Writer
}

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
