# Architecture

The idea behind atwhy is to have several interfaces, each one for a small purpose. These interfaces are then
implemented by concrete implementations.

It is possible to replace or mock each part of the application at any time. You can use atwhy as lib and just provide
your own implementations.

The interfaces are:
* `Loader` loads files from a given path.  
* `loader.TagFinder` reads the file and returns all lines which are part of a found tag. It Does not process the raw lines.  
* `TagFactories` convert the raw tags from the `TagFinder` and generates final Tags out of them.  
* `TemplateLoader` loads the templates from the `template` folder to pass them the generator.  
* `Generator` is responsible for postprocessing the tags and output the final file. which it just writes to the  
passed `Writer`.  
  
So the workflow is:  
Loader -> TagFinder = tagList []tag.Raw tagList -> TagProcessor -> TemplateLoader -> Generator -> Writer  
[core/atwhy.go:44]( core/atwhy.go )  
```go
type AtWhy struct {
	Loader         Loader
	Finder         loader.TagFinder
	TagFactories   []tag.Factory
	Generator      Generator
	TemplateLoader TemplateLoader

	projectPath       string
	projectPathPrefix string
	pageTemplate      *template.Template
}
```

