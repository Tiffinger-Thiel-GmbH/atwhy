package main

import (
	"strings"
)

type SpacePrefixCleaner struct{}

func (s SpacePrefixCleaner) Clean(in string) (string, error) {
	lines := strings.Split(in, "\n")

	for i := range lines {

		// Delete the first space after the special chars (if there is one)
		lines[i] = strings.TrimPrefix(lines[i], " ")
	}

	return strings.Join(lines, "\n"), nil
}
