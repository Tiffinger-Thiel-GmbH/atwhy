package finder

import (
	"bufio"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
)

type CommentConfig struct {
	LineComment []string
	BlockStart  []string
	BlockEnd    []string
}

// Finder implements the TagFinder interface in a language-generic way.
// It takes into account block comments (e.g. /* .... */) and line comments
// (e.g. // ...). You can pass alternative comment indicators to
// support other languages.
type Finder struct {
	// CommentConfig maps the filetype (e.g. ".go") to the matching CommentConfig.
	CommentConfig map[string]CommentConfig

	currentlyInBlockComment  bool
	currentLineIsLineComment bool

	// currentCommentLine is the current line cleaned up from the comment-indicators.
	// It is empty if the current line is no comment.
	currentCommentLine string

	// currentBlockIndex of the currently used block BlockCommentStart
	// The BlockCommentEnd with the same index closes the comment.
	currentBlockIndex int

	currentTag *tag.Raw

	// includeCode saves if a \@WHY CODE tag was found.
	// It has to be reset at a \@WHY CODE_END tag.
	includeCode bool
}

func (f *Finder) finishTag(res []tag.Raw) []tag.Raw {
	if f.currentTag != nil {

		f.currentTag.Value = f.currentTag.Value + f.currentCommentLine

		// Unescape \@ to @
		f.currentTag.Value = strings.ReplaceAll(f.currentTag.Value, "\\@", "@")
		res = append(res, *f.currentTag)
		f.currentTag = nil
	}

	return res
}

func (f *Finder) reset() {
	f.currentlyInBlockComment = false
	f.currentLineIsLineComment = false
	f.currentCommentLine = ""
	f.currentBlockIndex = -1
	f.currentTag = nil
	f.includeCode = false
}

func (f *Finder) Find(filename string, reader io.Reader) ([]tag.Raw, error) {
	f.reset()

	var res []tag.Raw
	var scan = bufio.NewScanner(reader)

	var lineNum = -1
	for scan.Scan() {
		lineNum++

		line := scan.Text()
		commentCFG, found := f.CommentConfig[filepath.Ext(filename)]
		if !found {
			commentCFG, found = f.CommentConfig["."]
			if !found {
				return res, nil
			}
		}
		f.findComment(commentCFG, line)

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
				}

				// Special tag CODE_END
				if newTag.Type == tag.TypeCodeEnd {
					f.includeCode = false
					f.currentCommentLine = ""
					res = f.finishTag(res)
					continue
				}

				// Finish the previous tag and start a new one.
				f.currentCommentLine = ""
				res = f.finishTag(res)
				newTag.Filename = filename
				newTag.Line = lineNum
				f.currentTag = newTag

				// Special tag LINK doesn't need any additional lines,
				// we can stop here.
				if newTag.Type == tag.TypeLink {
					res = f.finishTag(res)
					continue
				}

				continue
			}

			// If includeCode == false just add the current line to the value.
			if f.currentTag != nil && !f.includeCode {
				f.currentTag.Value = f.currentTag.Value + f.currentCommentLine + "\n"
				continue
			}
		}

		// If includeCode == true, add the whole line without trimming.
		if f.currentTag != nil && f.includeCode {
			f.currentTag.Value = f.currentTag.Value + line + "\n"
			continue
		}

		// For empty comment lines, just add newlines.
		if f.currentTag != nil && f.currentCommentLine == "" && (f.currentlyInBlockComment || f.currentLineIsLineComment) {
			f.currentTag.Value = f.currentTag.Value + "\n"
			continue
		}

		f.currentCommentLine = ""
	}

	// Finish the last tag.
	f.currentCommentLine = ""
	res = f.finishTag(res)

	return res, nil
}

// findComment and sets the struct-variables
// currentCommentLine, currentLineIsLineComment, currentBlockIndex, currentlyInBlockComment
// accordingly.
func (f *Finder) findComment(cfg CommentConfig, line string) {
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
		for _, lineCommentStart := range cfg.LineComment {
			if !f.currentlyInBlockComment && strings.HasPrefix(trimmedLine, lineCommentStart) {
				f.currentLineIsLineComment = true
				f.currentCommentLine = strings.TrimLeft(trimmedLine, lineCommentStart)
				return
			}
		}

		// Then check if it is in a block comment.
		for blockIndex, blockCommentStart := range cfg.BlockStart {
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

		for _, blockCommentEnd := range cfg.BlockEnd {

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
// Matches @ or \@ with any following postfix:
//
//	DOC any_name
//	DOC CODE any_name
//	DOC CODE_END
//	DOC LINK any_name
//
// @WHY readme_tags-2_rules
// The placeholder_names must follow these rules:
// First char: only a-z (lowercase)
// Rest:
//   - only a-z (lowercase)
//   - `-`
//   - `_`
//   - 0-9
//
// Examles:
//   - any_tag_name
//   - supertag
//   - super-tag
var anyTagRegex = regexp.MustCompile(`([\\]?)@WHY( ([A-Z_]+))?( ([a-z]+[a-z-_0-9]*))?`)

// findTag using the anyTagRegex.
// It already pre-fills the first comment line if a new one was found.
func (f *Finder) findTag() *tag.Raw {
	matches := anyTagRegex.FindAllStringSubmatch(f.currentCommentLine, 1)
	if matches == nil {
		return nil
	}

	// Per line, we only need one match.
	match := matches[0]

	// Ignore escaped \@WHY
	if match[1] != "" {
		return nil
	}

	newTag := tag.Raw{
		Type:        tag.Type(match[3]),
		Placeholder: match[5],
		Value:       f.currentCommentLine + "\n",
	}

	// If none was given, it is a DOC.
	if newTag.Type == "" {
		newTag.Type = tag.TypeDoc
	}

	return &newTag
}
