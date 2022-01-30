# atwhy [![test](https://github.com/Tiffinger-Thiel-GmbH/atwhy/actions/workflows/test.yaml/badge.svg)](https://github.com/Tiffinger-Thiel-GmbH/atwhy/actions/workflows/test.yaml)

## What is atwhy

atwhy can be used to generate a documentation out of comments in the code.  
That way you can for example describe all available options in the same file  
where they are coded. A developer therefore doesn't have to know exactly where  
the information has to be documented because it is just in the same file.

--> __Single source of truth__ also for documentation!

Markdown templates are used to just group the tags together into files and 
add some non-code specific docu.

The idea of athwy was born during a company-hackathon and since then evolved to a first fully usable 
preview version.

__Although most things are in a stable state, there may be small breaking changes until
awhy reaches v1.0.0.__

## Example

As __atwhy__ itself uses __atwhy__, you can just 
* look at the [templates](templates) folder of this project.
* and search for `@WHY` in the whole project.

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
atwhy serve  
```  
For more information run `atwhy serve --help`

### Templates

The templates are by default inside the [templates](templates) folder of your project.  
Each template results in one .md file.
So if you have the file `templates/README.tpl.md` it will be generated to `README.md`.
If you have the file `templates/doc/Usage.tpl.md` it will be generated to `doc/Usage.md`.

The templates should be markdown files with a yaml header for metadata.  
  
You can access a tag called `@WHY example_tag` using  
 ```text  
 # Example  
 {{ .Tag.example_tag }}  
 ```  
  
Note: This uses the Go templating engine.  
Therefor you can use the [Go templating syntax](https://learn.hashicorp.com/tutorials/nomad/go-template-syntax?in=nomad/templates).

__Possible template values are:__  
* Any Tag from the project: `{{ .Tag.example_tag }}`  
* Current Datetime: `{{ .Now }}`  
* Metadata from the yaml header: `{{ .Meta.Title }}`  
* Conversion of links to project-files (also in serve-mode): `{{ .Project "my/file/in/the/project.go" }}`  
  You need to use that if you want to generate links to actual files in your project.  
  This can also be used for pictures: `![aPicture]({{ .Project "path/to/the/picture.jpg" }})`  

__What if `{{` or `}}` is needed in the documentation?__  
You can wrap them with `{{.Escape "..."}}`.  
E.g.: `{{ .Escape "\"{{\"  and  \"}}\"" }}`  
Results in this markdown text: `"{{" and "}}"`  
  
__Note:__ You need to escape `"` with `\"`.  
  
(The official Go-Template way `{{ "{{ -- }}" }}` doesn't work in all cases with atwhy. `.Escape` works always.)

#### Header

Each template may have a yaml Header.  
Example with all possible fields:  
```markdown  
---  
# Some metadata which may be used for the generation.  
meta:  
  # The title is used for the served html to e.g. generate a menu and add page titles.  
  title: Readme # default: the template filename  
  
# Additional configuration for the `atwhy serve` command.  
server:  
  index: true # default: false  
---  
# Your Markdown starts here  
  
## Foo  
bar  
```  
(Note: VSCode supports the header automatically.)  

### Tags

Tags are the heart of __atwhy__.  
Basically you can add them in any comment of any file and then reference them
in any of the templates.
(Currently only `//` and `/*  */` comments are supported, but this will change.)

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

* \>= Go 1.16

### Build

Run `go build .`  

---
Generated: __30 Jan 22 14:43 +0100__
