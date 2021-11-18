package main

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type Finder struct {
	BlockCommentStarts []string
	BlockCommentEnds   []string
	LineCommentStarts  []string

	currentlyInBlockComment  bool
	currentLineIsLineComment bool
	currentCommentLine       string
	currentBlockIndex        int

	currentTag *tag.Raw
}

func (f *Finder) finishTag(res []tag.Raw) []tag.Raw {
	if f.currentTag != nil {
		f.currentTag.Value = f.currentTag.Value + f.currentCommentLine + "\n"
		res = append(res, *f.currentTag)
		f.currentTag = nil
	}

	return res
}

/** @README
  *

/* fthzfhgfth

dgf
dfg
df
g
dsgfdg*/
func (f *Finder) Find(filename string, reader io.Reader) ([]tag.Raw, error) {
	var res []tag.Raw
	var scan = bufio.NewScanner(reader)

	var lineNum int = -1
	for scan.Scan() {
		lineNum++

		line := scan.Text()
		f.findComment(line)

		if !f.currentlyInBlockComment &&
			!f.currentLineIsLineComment {
			res = f.finishTag(res)
			continue
		}

		if f.currentCommentLine != "" {
			newTag := f.findTag()
			if newTag != nil {
				f.currentCommentLine = ""
				res = f.finishTag(res)
				newTag.Filename = filename
				newTag.Line = lineNum
				f.currentTag = newTag
				continue
			}

			if f.currentTag != nil {
				f.currentTag.Value = f.currentTag.Value + f.currentCommentLine + "\n"
				continue
			}
		}

		if f.currentTag != nil && f.currentCommentLine == "" && (f.currentlyInBlockComment || f.currentLineIsLineComment) {
			f.currentTag.Value = f.currentTag.Value + "\n"
			continue
		}

		f.currentCommentLine = ""
	}

	f.currentCommentLine = ""
	res = f.finishTag(res)

	return res, nil
}

func (f *Finder) findComment(line string) {
	f.currentCommentLine = ""
	f.currentLineIsLineComment = false

	var commentStartedInThisLine bool
	trimmedLine := strings.TrimLeft(line, " ")
	if !f.currentlyInBlockComment {

		for _, lineCommentStart := range f.LineCommentStarts {
			if strings.HasPrefix(trimmedLine, lineCommentStart) {
				f.currentLineIsLineComment = true
				f.currentCommentLine = strings.TrimLeft(trimmedLine, lineCommentStart)
				return
			}
		}

		for blockIndex, blockCommentStart := range f.BlockCommentStarts {
			if strings.HasPrefix(trimmedLine, blockCommentStart) {
				f.currentBlockIndex = blockIndex
				f.currentlyInBlockComment = true
				f.currentCommentLine = strings.TrimLeft(trimmedLine, blockCommentStart)
				commentStartedInThisLine = true
				break
			}
		}
	}

	if f.currentlyInBlockComment {
		linePart := trimmedLine
		if commentStartedInThisLine {
			linePart = f.currentCommentLine
		}

		for _, blockCommentEnd := range f.BlockCommentEnds {

			if foundIndex := strings.Index(linePart, blockCommentEnd); foundIndex > -1 {
				f.currentBlockIndex = -1
				f.currentlyInBlockComment = false
				f.currentCommentLine = linePart[:foundIndex]
				return
			}
		}

		f.currentBlockIndex = -1
		f.currentCommentLine = linePart
	}
}

var anyTagRegex = regexp.MustCompile("@[A-Za-z]+")

func (f *Finder) findTag() *tag.Raw {
	if tagType := anyTagRegex.FindString(f.currentCommentLine); tagType != "" {
		return &tag.Raw{
			Type:  tag.Type(tagType[1:]),
			Value: f.currentCommentLine + "\n",
		}
	}

	return nil
}

// func (f Finder) findTags(fileName string, text []string) []tag.Raw {
// 	var foundTagLine bool
// 	var taggys []string
// 	var tagLine int
// 	var tagValue string
// 	var tagType string
// 	var finalTags []tag.Raw

// 	for index, eachLn := range text {

// 		if f.findTagLines(eachLn) != "" {
// 			foundTagLine = true
// 			tagType = f.findTagLines(eachLn)
// 			tagLine = index + 1
// 		}

// 		if foundTagLine {

// 			if len(strings.TrimSpace(eachLn)) == 0 {

// 				tagValue = strings.Join(taggys, " \n ")
// 				finalTags = append(finalTags, tag.Raw{Type: tag.Type(tagType), Filename: fileName, Line: tagLine, Value: tagValue})
// 				taggys = nil
// 				foundTagLine = false

// 			} else {

// 				taggys = append(taggys, eachLn)
// 			}

// 		}
// 	}
// 	return finalTags
// }

// /*
// @README
// hallo
// reame
// test
// */

// // @FILELINK
// // hallo

// TODO: pass if we are currently inside a block comment (as they go over several lines)
func (f Finder) findTagInLine(line string) string {
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
