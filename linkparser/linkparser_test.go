package linkparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestLinks(t *testing.T) {
	info, err := ioutil.ReadDir("tests")
	if err != nil {
		t.Fatalf("Can't find directory 'tests'. Err - %s", err)
	}
	for _, f := range info {
		if !f.IsDir() {
			file, err := os.Open("tests/" + f.Name())
			if err != nil {
				t.Fatal(err)
			}
			links, err := Links(file)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(links)
		}
	}
}
