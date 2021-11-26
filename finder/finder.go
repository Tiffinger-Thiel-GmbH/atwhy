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

	// includeCode saves if a \@CODE tag was found.
	// It has to be reset at a \@CODE_END tag.
	includeCode bool
}

func (f *Finder) finishTag(res []tag.Raw) []tag.Raw {
	if f.currentTag != nil {
		f.currentTag.Value = f.currentTag.Value + f.currentCommentLine + "\n"

		// Unescape \@ to @
		f.currentTag.Value = strings.ReplaceAll(f.currentTag.Value, "\\@", "@")
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

		// Finish the current tag if there is no more comment line (or includeCode).
		if !f.currentlyInBlockComment &&
			!f.currentLineIsLineComment &&
			!f.includeCode {
			res = f.finishTag(res)
			continue
		}

		if f.currentCommentLine != "" {
			newTag := f.findTag()
			if newTag != nil {
				// Special tag CODE
				if newTag.Type == tag.TypeCode {
					f.includeCode = true
					continue
				}

				// Special tag CODE_END
				if newTag.Type == tag.TypeCodeEnd {
					f.includeCode = false
					continue
				}

				// Finish the previous tag and start a new one.
				f.currentCommentLine = ""
				res = f.finishTag(res)
				newTag.Filename = filename
				newTag.Line = lineNum
				f.currentTag = newTag

				continue
			}

			// Just add the current line to the value.
			if f.currentTag != nil {
				f.currentTag.Value = f.currentTag.Value + f.currentCommentLine + "\n"
				continue
			}
		}

		// For empty comment lines, just add newlines.
		if f.currentTag != nil && f.currentCommentLine == "" && (f.currentlyInBlockComment || f.currentLineIsLineComment) {
			f.currentTag.Value = f.currentTag.Value + "\n"
			continue
		}

		// If no longer in comment but still includeCode, add the whole line as code.
		if f.currentTag != nil && f.includeCode {
			f.currentTag.Code = f.currentTag.Code + line + "\n"
			continue
		}

		f.currentCommentLine = ""
	}

	// Finish the last tag.
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

		// First check if it is a One-Line comment. (e.g. //)
		for _, lineCommentStart := range f.LineCommentStarts {
			if strings.HasPrefix(trimmedLine, lineCommentStart) {
				f.currentLineIsLineComment = true
				f.currentCommentLine = strings.TrimLeft(trimmedLine, lineCommentStart)
				return
			}
		}

		// Then check if it is in a block comment.
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

	// Try to find the end of the block comment.
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

// anyTagRegex matches all tags include a possible \ which is then checked as it escapes the tag.
var anyTagRegex = regexp.MustCompile(`[\\]?@[A-Za-z_]+`)

func (f *Finder) findTag() *tag.Raw {
	if tagType := anyTagRegex.FindString(f.currentCommentLine); tagType != "" && tagType[0] != '\\' {
		return &tag.Raw{
			Type:  tag.Type(tagType[1:]),
			Value: f.currentCommentLine + "\n",
		}
	}

	return nil
}
