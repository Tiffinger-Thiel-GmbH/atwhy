# README

## What is CrazyDoc

CrazyDoc can be used to generate a documentation out of comments in the code.  
That way you can for example describe all available options in the same file  
where they are coded. A developer therefore doesn't have to know exactly where  
the information has to be documented because it is just in the same file.

The same applies to architectural decisions, which can be documented, where its  
actually done.  
--> __Single source of truth__ also for documentation!

## Usage

{{ .Tag.readme_usage }}

## Ignore

* You can pass something like `--ext ".go,.js,.ts"` to only process specific files.
* You can create a `.crazydocignore` file which just follows the `.gitignore` syntax.  
  (If you find an inconsistency with the git-handling, please report it [here](https://github.com/aligator/NoGo/issues).)

## Tags

{{ .Tag.readme_tags }}

The tag names must follow the following rules:
{{ .Tag.readme_tags_rules }}

The tags are terminated by 
* another tag
* empty line
* Exception: `@DOC CODE` is terminated by `@DOC CODE_END` and not by empty lines.

## Distribute

### Prerequisites

* Go 1.17

### Build

Run `go build .`  