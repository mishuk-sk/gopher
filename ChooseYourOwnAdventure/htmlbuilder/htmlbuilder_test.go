package htmlbuilder

import (
	"fmt"
	"testing"

	"github.com/mishuk-sk/gopher/ChooseYourOwnAdventure/htmlbuilder/storyparser"
)

func TestGetPage(t *testing.T) {
	arc := storyparser.Arc{
		Label:      "info",
		Title:      "first try",
		Paragraphs: []string{"Paragraph 1", "Paragraph 2", "Paragraph 3"},
		Options: []storyparser.Option{
			{
				Text: "First Option",
				Link: "1",
			},
			{
				Text: "Second Option",
				Link: "2",
			},
		},
	}
	fmt.Println(string(GetPage(arc)))
}
