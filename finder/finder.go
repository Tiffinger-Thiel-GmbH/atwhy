package finder

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
	defer func() {
		// Always cut the first space because usually comments have a space after the comment sign.
		f.currentCommentLine = strings.TrimPrefix(f.currentCommentLine, " ")
	}()

	f.currentCommentLine = ""
	f.currentLineIsLineComment = false

	var commentStartedInThisLine bool
	trimmedLine := strings.TrimLeft(line, " \t")
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
