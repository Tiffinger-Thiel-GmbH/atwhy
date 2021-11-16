package main

import (
	"strings"
)

type SlashStarCleaner struct{}

func (s SlashStarCleaner) Clean(in string) (string, error) {
	lines := strings.Split(in, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimLeft(strings.TrimLeft(line, " "), "*/")
	}

	return strings.Join(lines, "\n"), nil
}
