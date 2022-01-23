# Architecture

The idea behind atwhy is to have several interfaces, each one for a small purpose. These interfaces are then
implemented by concrete implementations.

It is possible to replace or mock each part of the application at any time. You can use atwhy as lib and just provide
your own implementations.

The interfaces are:

* `Loader` loads files from a given path.
* `TagFinder` reads the file and returns all lines which are part of a found tag. It Does not process the raw lines.
* `TagProcessor` processes the raw data from the `TagFinder` and generates Tags out of them. It may also clean
  comment-chars and spaces and combine some tags.
* `Generator` is responsible for postprocessing the tags and output the final file. which it just writes to the
  passed `Writer`.

So the workflow is:
Loader -> TagFinder = tagList []tag.Raw tagList -> TagProcessor -> Generator -> Writer

[../core/atwhy.go:34](../core/atwhy.go)  
```go
type AtWhy struct {
	Loader         Loader
	Finder         loader.TagFinder
	Processor      TagProcessor
	Generator      Generator
	TemplateLoader TemplateLoader

	projectPath  string
	pageTemplate *template.Template
}
```

