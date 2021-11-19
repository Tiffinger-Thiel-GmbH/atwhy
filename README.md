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
Just run `crazydoc [OPTIONS]... [PROJECT_ROOT]`.  
To get all possible file extensions just run `crazydoc -help`  

## Distribute
### Prerequisites  
* Go 1.17  
### Build  
Run `go build .`  

