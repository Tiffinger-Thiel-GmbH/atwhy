package main

//go:generate go run . -out README.md -tags README
import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/finder"
	"gitlab.com/tiffinger-thiel/crazydoc/generator"
	"gitlab.com/tiffinger-thiel/crazydoc/loader"
	"gitlab.com/tiffinger-thiel/crazydoc/processor"

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

// @README 40
// Distribute
// # Prerequisites
// * Go 1.17
//
// # Build
// Run `go build .`

type Loader interface {
	Load(dir string, finder loader.TagFinder) (allTags []tag.Raw, err error)
}

type TagProcessor interface {
	Process(tags []tag.Raw) ([]tag.Tag, error)
}

type Generator interface {
	Generate(tags []tag.Tag, writer io.Writer) error
}

func New(fileExtensions []string, tagsToExport []string, outputFile string) CrazyDoc {
	var gen Generator

	outputFileExtension := filepath.Ext(outputFile)

	switch outputFileExtension {
	case ".md", "":
		gen = &generator.Markdown{
			TagsToExport: tagsToExport,
		}
	case ".html":
		gen = &generator.HTML{
			Markdown: generator.Markdown{
				TagsToExport: tagsToExport,
			},
		}
	}

	writer := os.Stdout

	crazyDoc := CrazyDoc{
		Finder: &finder.Finder{
			BlockCommentStarts: []string{"/*"},
			BlockCommentEnds:   []string{"*/"},
			LineCommentStarts:  []string{"//"},
		},
		Loader: loader.File{
			FS:             afero.NewOsFs(),
			FileExtensions: fileExtensions,
		},
		Processor: processor.Processor{
			TagFactories: []tag.Factory{
				tag.Why,
				tag.Readme,
				tag.FileLink,
			},
		},
		Generator: gen,
		Writer:    writer,
	}
	return crazyDoc
}

func ParseCmd() (fileExtensions []string, tagsToExport []string, outputFile string, inputPath string, host string) {
	// @README 20
	// Usage
	// Just run `crazydoc [OPTIONS]... [PROJECT_ROOT]`.
	// To get all possible options just run `crazydoc -help`
	//
	// In development, use `go run . [OPTIONS]... [PROJECT_ROOT]` instead.

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS]... [PROJECT_ROOT]\n", os.Args[0])
		flag.PrintDefaults()
	}
	extVar := flag.String("ext", ".go,.js,.ts,.jsx,.tsx", "comma separated list of file extensions to search for")
	tagTypes := flag.String("tags", "README,WHY", "comma separated list tag types that should be exported")
	hostVar := flag.String("host", "", "serves generated html file to given host (e.g. localhost:4000) \n-out param will be ignored")
	outputFileVar := flag.String("out", "", "ouptut file \nshould be a .md or .html file")
	flag.Parse()

	inputPath = flag.Arg(0)
	fileExtensions = strings.Split(*extVar, ",")
	tagsToExport = strings.Split(*tagTypes, ",")
	if inputPath == "" {
		inputPath = "."
	}

	return fileExtensions, tagsToExport, *outputFileVar, inputPath, *hostVar
}

func main() {
	fileExtensions, tagsToExport, outputFile, inputPath, host := ParseCmd()

	if host == "" {
		crazyDoc := New(fileExtensions, tagsToExport, outputFile)

		if outputFile != "" {
			file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			crazyDoc.Writer = file
		}

		if err := crazyDoc.Run(inputPath); err != nil {
			panic(err)
		}
		return
	}

	fs := http.FileServer(http.Dir(inputPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fs.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

		crazyDoc := New(fileExtensions, tagsToExport, outputFile)
		crazyDoc.Generator = &generator.HTML{
			Markdown: generator.Markdown{
				TagsToExport: tagsToExport,
			},
		}
		crazyDoc.Writer = w

		if err := crazyDoc.Run(inputPath); err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println("Starting server")
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Fatal(err)
	}

}
