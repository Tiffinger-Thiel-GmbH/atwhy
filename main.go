//go:generate go run . --ext .go --templates README README.md

package main

import "github.com/Tiffinger-Thiel-GmbH/atwhy/cmd"

func main() {
	cmd.Execute()
}
