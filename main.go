//go:generate go run . --templates-folder docTemplates --ext .go --templates README README.md

package main

import "gitlab.com/tiffinger-thiel/crazydoc/cmd"

func main() {
	cmd.Execute()
}
