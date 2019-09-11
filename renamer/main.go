package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var reg *regexp.Regexp
var replaceString = "$2-$1$3"

func main() {
	p := flag.String("path", "./sample", "Used to provide starting point for renaming")
	regexString := flag.String("re", "(.+?)_([0-9]+)(.+)", "Used to provide regular expression for files to replace")
	replStr := flag.String("str", "$2-$1$3", "New order of matching groups from regexp")
	flag.Parse()
	reg = regexp.MustCompile(*regexString)
	replaceString = *replStr
	filepath.Walk(*p, walkFunc(rename))
}

func walkFunc(renameFunc func(string, string) error) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Visited with error: %v\n", err)
		}
		if !info.IsDir() {
			if err := renameFunc(path, info.Name()); err != nil {
				log.Println(err)
			}
		}
		return nil
	}
}

func rename(path, name string) error {
	//FIXME
	// if files are renaming to the same name - only last one will be presented
	if !reg.MatchString(name) {
		return fmt.Errorf("filename <%s> doesn't match regular expression <%s>", name, reg.String())
	}
	p := strings.TrimSuffix(path, name)
	newName := reg.ReplaceAllString(name, replaceString)
	return os.Rename(path, filepath.Join(p, newName))
}
