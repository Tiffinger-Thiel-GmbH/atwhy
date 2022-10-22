---
meta:
  title: Readme 
server:
  index: true
---
# atwhy [![test](https://github.com/Tiffinger-Thiel-GmbH/atwhy/actions/workflows/test.yaml/badge.svg)](https://github.com/Tiffinger-Thiel-GmbH/atwhy/actions/workflows/test.yaml) [![codecov](https://codecov.io/gh/Tiffinger-Thiel-GmbH/atwhy/branch/main/graph/badge.svg?token=JSN8ANHSNA)](https://codecov.io/gh/Tiffinger-Thiel-GmbH/atwhy)

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
atwhy reaches v1.0.0.__

## Example

As __atwhy__ itself uses __atwhy__ very extensively, you can just
* look at the [templates]({{ .Project "templates" }}) folder of this project.
* and [search for `@WHY`](https://github.com/Tiffinger-Thiel-GmbH/atwhy/search?q=%5C%40WHY&type=) in the whole project.

## Installation

You have several options to install atwhy:
* Just use docker to run a minimal image (It is multi-arch so it works on x64 and arm):  
  `docker run --rm -i -p 4444:4444 -u $UID:$GID -v $PWD:/project ghcr.io/tiffinger-thiel-gmbh/atwhy atwhy`  
  You may add an alias to your shell for this.
* [Install Go](https://go.dev/dl/) and run `go install github.com/Tiffinger-Thiel-GmbH/atwhy@latest`.  
  You may need to restart after installing Go to have the PATH setup correctly.
* [Download a matching binary from the releases](https://github.com/Tiffinger-Thiel-GmbH/atwhy/releases)
  and put somewhere in your PATH.  
  Note that they are currently not signed -> MacOS and Windows may not like that...

## Usage

### Command

{{ .Tag.readme_usage }}

{{ .Tag.readme_usage_serve }}

### Templates

The templates are by default inside the [templates]({{ .Project "templates" }}) folder of your project.  
Each template results in one .md file.
So if you have the file `templates/README.tpl.md` it will be generated to `README.md`.
If you have the file `templates/doc/Usage.tpl.md` it will be generated to `doc/Usage.md`.

{{ .Tag.doc_template }}

{{ .Tag.doc_template_possible_tags }}  

{{ .Tag.doc_template_escape_tag }}

#### Header

{{ .Tag.doc_template_header_1 }}  

### Tags

Tags are the heart of __atwhy__.  
Basically you can add them in any comment of any file and then reference them
in any of the templates.

{{ .Tag.readme_tags }}

{{ .Tag.readme_tags_rules }}

The tags are terminated by

* another tag
* empty line
* Exception: `@WHY CODE` is terminated by `@WHY CODE_END` and not by empty lines.

### Comments

You can specify the type of comments for each type of file.
For this you may use the `--comment` flag.

{{ .Tag.readme_comments }}

The following are the default, built-in rules:
{{ .Tag.readme_comments_builtin }}

### Ignore

* You can pass something like `--ext ".go,.js,.ts"` to only process specific files.
* You can create a `.atwhyignore` file which just follows the `.gitignore` syntax.  
  (If you find an inconsistency with the git-handling, please report it 
  [here](https://github.com/aligator/NoGo/issues).)

## Distribute

### Prerequisites

* \>= Go 1.18

### Build

Run `go build .`  

---
This README was last updated on: __{{ .Now }}__
