package storyparser

import (
	"encoding/json"
	"fmt"
)

// Arc defines whole story arc with title, paragraphs,
//options and arc uniq label
type Arc struct {
	Label      string
	Title      string
	Paragraphs []string
	Options    []Option
}

// Option defines choices for next steps in story
type Option struct {
	Text string
	Link string
}

// ParseStory parses correct json file into slice of Arc
func ParseStory(data []byte) ([]Arc, error) {
	var story []Arc
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, fmt.Errorf("Error unmarshalling json to story object. %s", err)
	}
	return story, nil
}
