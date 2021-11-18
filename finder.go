package main

import (
	"bufio"
	"io"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type Finder struct {
	BlockCommentStarts []string
	BlockCommentEnds   []string
	LineCommentStarts  []string
}

func (f Finder) Find(filename string, reader io.Reader) ([]tag.Raw, error) {
	var scan = bufio.NewScanner(reader)
	scan.Split(bufio.ScanLines)
	var text []string
	for scan.Scan() {
		text = append(text, scan.Text())
	}

	var tags = f.findTags(filename, text)

	return tags, nil
}

func (f Finder) findTags(fileName string, text []string) []tag.Raw {
	var foundTagLine bool
	var taggys []string
	var tagLine int
	var tagValue string
	var tagType string
	var finalTags []tag.Raw

	for index, eachLn := range text {

		if f.findTagLines(eachLn) != "" {
			foundTagLine = true
			tagType = f.findTagLines(eachLn)
			tagLine = index + 1
		}

		if foundTagLine {

			if len(strings.TrimSpace(eachLn)) == 0 {

				tagValue = strings.Join(taggys, " \n ")
				finalTags = append(finalTags, tag.Raw{Type: tag.Type(tagType), Filename: fileName, Line: tagLine, Value: tagValue})
				taggys = nil
				foundTagLine = false

			} else {

				taggys = append(taggys, eachLn)
			}

		}
	}
	return finalTags
}

// TODO: pass if we are currently inside a block comment (as they go over several lines)
func (f Finder) findTagLines(line string) string {
	// TODO: remember if it is currently a comment

	var foundPossibleTag bool
	var tagName string

	// TODO: If we are not in a comment:
	// Check if the line starts with a LineComment - ignore spaces

	for charIndex, char := range line {
		// TODO: If not already in comment.
		// Check if the current position opens one (note: it may be multi-char)

		// If you are now in a comment, continue:

		// TODO: If not already in comment.
		// Check if the current position closes a comment (note: it may be multi-char)

		// If you are still in a comment, continue:

		if foundPossibleTag {
			if char == ' ' || charIndex == len(line) {
				return tagName
			}

			tagName += string(char)
			continue
		}

		if char == '@' {
			foundPossibleTag = true
			continue
		}
	}

	return ""
}
