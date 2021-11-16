package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type Finder struct {
}

func (f Finder) Find(filename string, reader io.Reader) ([]tag.Raw, error) {

	filename = "test.js"

	return []tag.Raw{}, nil
}

func Example() {
	const src = `
	// type TagFinder interface {
	// Find(filename string, reader io.Reader) ([]Tag, error)
	// Find(a int, b int ) (int, error)
	// SaveByTag()
	// @Test hier steht mein tag
	// scan()
	// type TagFinder interface {
	// Find(filename string, reader io.Reader) ([]Tag, error)
	// Find(a int, b int ) (int, error)
	// SaveByTag()

	// @README hier steht mein tag
	 findTag()
	
	 saveByTag()
`

   const filename = "test.js"

	var scan = bufio.NewScanner(strings.NewReader(src))
	scan.Split(bufio.ScanLines)
	var text []string
    for scan.Scan() {
        text = append(text, scan.Text())
    }

	var tags = findTags(text)
	
	fmt.Println(tags)
}


func findTags(text []string) []tag.Raw {
	var foundTagLine bool;
	var taggys []string;
	var tagLine int; 
	var tagValue string;
	var tagType string;
	var finalTags []tag.Raw

	for index, eachLn := range text {

		if(findTagLines(eachLn) != ""){
			foundTagLine = true
			tagType = findTagLines(eachLn)
			tagLine = index;

		}

		if(foundTagLine){

			
			if(len(strings.TrimSpace(eachLn)) == 0 ){
			
			tagValue = strings.Join(taggys, " ")
			finalTags = append(finalTags, tag.Raw {Type: tag.Type(tagType), Filename: "af", Line: tagLine, Value: tagValue})
			
			
			}

			taggys = append(taggys, eachLn)

		}
		
	}
	
	
	return finalTags
}

func findTagLines(eachLn string) string {
	var foundPossibleTag bool = false 
	var tagName string = ""

	for _, rune := range eachLn {
			
		if foundPossibleTag {
			if rune == ' '{
				return tagName
			}

			tagName += string(rune)
			continue 
		}

		if rune == '@'{
			foundPossibleTag  = true
			continue
		}
	}

	return ""
}