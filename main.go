package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

// @README 10
// What is CrazyDoc
// CrazyDoc can be used to generate a documentation out of comments in the code.
// That way you can for example describe all available options in the same file
// where they are coded. A developer therefore doesn't have to know exactly where
// the information has to be documented because it is just in the same file.
//
// The same applies to architectural decisions, which can be documented, where its
// actually done.
// --> __Single source of truth__ also for documentation!

// @README 20
// Distribute
// # Prerequisites
// * Go 1.17
//
// # Build
// Run `go build .`
//

type Loader interface {
	Load(dir string, finder TagFinder) (allTags []tag.Raw, err error)
}

type TagFinder interface {
	Find(filename string, reader io.Reader) (tags []tag.Raw, err error)
	// SaveByTag()
	// scan()
	// findTag()
	// saveByTag()
}

type TagProcessor interface {
	Process(tags []tag.Raw) ([]tag.Tag, error)
}

type Generator interface {
	Generate(tags []tag.Tag, writer io.Writer) error
}

func ParseCmd() (fileExtensions []string, tagsToExport []string, outputFile string, inputPath string) {
	// @README 20
	// Usage
	// Just run `crazydoc [OPTIONS]... [PROJECT_ROOT]`.
	// To get all possible file extensions just run `crazydoc -help`

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS]... [PROJECT_ROOT]\n", os.Args[0])
		flag.PrintDefaults()
	}
	extVar := flag.String("ext", ".go,.js,.ts,.jsx,.tsx", "comma separated list of file extensions to search for")
	tagTypes := flag.String("tags", "WHY,README", "comma separated list tag types that should be exported")
	outputFileVar := flag.String("out", "", "ouptut file \nshould be a .md or .html file")
	flag.Parse()

	inputPath = flag.Arg(0)
	fileExtensions = strings.Split(*extVar, ",")
	tagsToExport = strings.Split(*tagTypes, ",")
	if inputPath == "" {
		inputPath = "."
	}

	return fileExtensions, tagsToExport, *outputFileVar, inputPath
}

func main() {
	fileExtensions, tagsToExport, outputFile, inputPath := ParseCmd()

	var finder TagFinder = &Finder{
		BlockCommentStarts: []string{"/*"},
		BlockCommentEnds:   []string{"*/"},
		LineCommentStarts:  []string{"//"},
	}
	var loader Loader = FileLoader{
		FS:             afero.NewOsFs(),
		FileExtensions: fileExtensions,
	}
	var processor TagProcessor = Processor{
		cleaners: []Cleaner{
			SlashStarCleaner{},
		},
		tagFactories: []tag.Factory{
			tag.Why,
			tag.Readme,
			tag.FileLink,
		},
	}

	var generator Generator

	outputFileExtension := filepath.Ext(outputFile)

	switch outputFileExtension {
	case ".md", "":
		generator = MarkdownGenerator{
			TagsToExport: tagsToExport,
		}
	case ".html":
		generator = HTMLGenerator{
			MarkdownGenerator{
				TagsToExport: tagsToExport,
			},
		}
	}

	writer := os.Stdout
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer = file
	}

	crazyDoc := CrazyDoc{
		Finder:    finder,
		Loader:    loader,
		Processor: processor,
		Generator: generator,
		Writer:    writer,
	}

	if err := crazyDoc.Run(inputPath); err != nil {
		panic(err)
	}
}
