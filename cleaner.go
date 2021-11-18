package main

import (
	"strings"
)

type SlashStarCleaner struct{}

func (s SlashStarCleaner) Clean(in string) (string, error) {
	lines := strings.Split(in, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimLeft(strings.TrimLeft(line, " "), "*/")

		// Delete the first space after the special chars (if there is one)
		lines[i] = strings.TrimPrefix(lines[i], " ")
	}

	return strings.Join(lines, "\n"), nil
}
