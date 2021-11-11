package main

import "fmt"

type FileLoader interface {
}

type TagFinder interface {
}

type CommentCleaner interface {
}

type Generator interface {
}

type Writer interface {
}

func main() {
	fmt.Println("Hello World")
}
