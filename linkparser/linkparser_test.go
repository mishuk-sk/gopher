package linkparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/net/html"
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
			links, err := Parse(file)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(links)
		}
	}
}

type Test struct {
	input    string
	expected string
}

func TestGetText(t *testing.T) {
	tests := []Test{
		{
			input: `some text
			and more`,
			expected: `some text
			and more`,
		},
		{
			input:    `<h1>Some text</h1><h2>more text</h2>`,
			expected: `Some textmore text`,
		},
		{
			input:    `<h1>Some text<h2>more text</h2></h1>`,
			expected: `Some textmore text`,
		},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		doc, err := html.Parse(r)
		if err != nil {
			t.Fatalf("Can't parse test html. Err - %s", err)
		}
		buf := strings.Builder{}
		getText(doc, &buf)
		if buf.String() != test.expected {
			t.Errorf("Expected %s to equal %s", buf.String(), test.expected)
		}
	}
}

func TestGetNodes(t *testing.T) {
	tests := []Test{
		{
			input:    `<body><a href="#"></a></body>`,
			expected: "1",
		},
		{
			input: `<body><a href="#">Text
			<a href="1"></a>
			</a></body>`,
			expected: "1",
		},
		{
			input:    `<a href="#">Text<a href="1">more<span>text</span></a></a>`,
			expected: "1",
		},
		{
			input: `<body>
			<a></a>
			<a></a>
			</body>`,
			expected: "2",
		},
		{
			input: `<body>
			<a></a>
			<a></a>
			</body`,
			expected: "2",
		},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		doc, err := html.Parse(r)
		if err != nil {
			t.Fatalf("Can't parse test html. Err - %s", err)
		}
		nodes := getNodes(doc)
		if strconv.Itoa(len(nodes)) != test.expected {
			t.Errorf("Expected %v to have length of %s", nodes, test.expected)
		}
	}
}
