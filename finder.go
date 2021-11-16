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


	var scan = bufio.NewScanner(reader)
	scan.Split(bufio.ScanLines)
	var text []string
    for scan.Scan() {
        text = append(text, scan.Text())
    }

	var tags = findTags(filename, text)

	return tags, nil;
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
	// @IMPORTANT
	// impi impi
	
	 saveByTag()
`

   const filename = "test.js"

   var finder = Finder{}

   tags, err :=finder.Find("test.js", strings.NewReader(src))
  
   if(err != nil ){
	fmt.Println("ERROR")
   }

   fmt.Println(tags)

	
}


func findTags(fileName string, text []string) []tag.Raw {
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
			tagLine = index + 1;
		}

		if(foundTagLine){

			
			if(len(strings.TrimSpace(eachLn)) == 0){
			
			tagValue = strings.Join(taggys, " \n ")
			finalTags = append(finalTags, tag.Raw {Type: tag.Type(tagType), Filename: fileName, Line: tagLine, Value: tagValue})
			taggys = nil
			foundTagLine = false
			
			} else{

				taggys = append(taggys, eachLn)
			}

		}
	}
	return finalTags
}

func findTagLines(eachLn string) string {
	var foundPossibleTag bool = false 
	var tagName string = ""

	for charIndex, rune := range eachLn {
			
		if foundPossibleTag {
			if rune == ' ' || charIndex == len(eachLn){
				return tagName
			}

			tagName += string(rune)
			continue 
		}

		if rune == '@' {
			foundPossibleTag  = true
			continue
		}
	}

	return ""
}