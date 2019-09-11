package storyparser

import (
	"encoding/json"
	"fmt"
)

// Arc defines whole story arc with title, paragraphs,
//options and arc uniq label
type Arc struct {
	Label      string   `json:"-"`
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options,omitempty"`
}

// Option defines choices for next steps in story
type Option struct {
	Text string `json:"text,omitempty"`
	Link string `json:"arc,omitempty"`
}

// ParseStory parses correct json file into slice of Arc
func ParseStory(data []byte) ([]Arc, error) {
	var story []Arc
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, fmt.Errorf("Error unmarshalling json to story object. %s", err)
	}
	return story, nil
}

// MappedStory parses json file and returns Label: Arc map pairs
func MappedStory(data []byte) (map[string]Arc, error) {
	story := make(map[string]Arc)
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, fmt.Errorf("Error unmarshalling json to story object. %s", err)
	}

	return story, nil
}
