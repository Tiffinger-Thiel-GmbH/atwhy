# atwhy

## What is atwhy

atwhy can be used to generate a documentation out of comments in the code.  
That way you can for example describe all available options in the same file  
where they are coded. A developer therefore doesn't have to know exactly where  
the information has to be documented because it is just in the same file.

The same applies to architectural decisions, which can be documented, where its  
actually done.  
--> __Single source of truth__ also for documentation!

## Installation

You have several options to install atwhy:
* Just use docker to run a minimal image (It is multi-arch so it works on x64 and arm):  
  `docker run --rm -i -p 4444:4444 -v $PWD:/project ghcr.io/tiffinger-thiel-gmbh/atwhy atwhy`  
  You may add an alias to your shell for this.
* [Install Go](https://go.dev/dl/) and run `go install github.com/Tiffinger-Thiel-GmbH/atwhy@latest`.  
  You may need to restart after installing Go to have the PATH setup correctly.
* [Download a matching binary from the releases](https://github.com/Tiffinger-Thiel-GmbH/atwhy/releases)
  and put somewhere in your PATH.  
  Note that they are currently not signed -> MacOS and Windows may not like that...

## Usage

### Command

Usage:  
Just run  
```bash  
atwhy --help  
```  
If nothing special is needed, just run the command without any arguments:  
```bash  
atwhy  
```  
It will use the default values and just work if a `templates` folder with some  
templates (e.g. `templates/README.tpl.md`) exists.

You can also serve the documentation on default host `localhost:4444` with:  
```bash  
atwhy serve --ext .go  
```  
For more information run `atwhy serve --help`

### Templates

The templates should be markdown files with a yaml header for metadata.  
  
You can access a tag called `@WHY example_tag` using  
 ```text  
 # Example  
 {{ .Tag.example_tag }}  
 ```  
  
Note: This uses the Go templating engine.  
Therefor you can use the [Go templating syntax](https://learn.hashicorp.com/tutorials/nomad/go-template-syntax?in=nomad/templates).

Possible template values are:  
* Any Tag from the project: `{{ .Tag.example_tag }}`  
* Current Datetime: `{{ .Now }}`

#### Header

Each template has a yaml Header with the following fields:  
```go

type Header struct {
	// Meta contains additional data which can be used by the generators.
	// It is also available inside the template for example with
	//  {{ .Meta.Title }}
	Meta MetaData `yaml:"meta"`

	Server ServerData `yaml:"server"`
}

type MetaData struct {
	// Title is for example used in the html generator to create the navigation buttons.
	// If not set, it will default to the template file-name (excluding .tpl.md)
	Title string `yaml:"title"`
}

type ServerData struct {
	// Index defines if this template should be used as "index.html".
	// Note that there can only be one page in each folder which is the index.
	Index bool `yaml:"index"`
}
```
  
The header is separated from the markdown by using a line with three `-` and a newline.  
Example:  
```md  
meta:  
 title: Readme  
---  
# Your Markdown  
  
## Foo  
bar  
```

### Tags

You can use `@WHY <placeholder_name>` and then use that placeholder in any template.  
There are also some special tags:  
* `@WHY LINK <placeholder_name>` can be used to just add a link to the file where the tag is in.  
* `@WHY CODE <placeholder_name>` can be used to reference any code.  
  It has to be closed by `@WHY CODE_END`

The placeholder_names must follow these rules:  
First char: only a-z (lowercase)  
Rest:  
 * only a-z (lowercase)  
 * `-`  
 * `_`  
 * 0-9  
  
Examles:  
 * any_tag_name  
 * supertag  
 * super-tag

The tags are terminated by

* another tag
* empty line
* Exception: `@WHY CODE` is terminated by `@WHY CODE_END` and not by empty lines.

### Ignore

* You can pass something like `--ext ".go,.js,.ts"` to only process specific files.
* You can create a `.atwhyignore` file which just follows the `.gitignore` syntax.  
  (If you find an inconsistency with the git-handling, please report it 
  [here](https://github.com/aligator/NoGo/issues).)

## Distribute

### Prerequisites

* Go 1.16

### Build

Run `go build .`  

---
Generated: __22 Jan 22 19:22 +0100__
