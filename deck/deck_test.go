package deck_test

import (
	"testing"

	"github.com/mishuk-sk/gopher/deck"
)

func TestNew(t *testing.T) {
	d := deck.New()
	if len(d) != 52 {
		t.Fatalf("Expected deck to be of 52 cards but got %d\n Deck: %v\n", len(d), d)
	}
	for i := deck.Spades; i <= deck.Hearts; i++ {
		for j := deck.Ace; j <= deck.King; j++ {
			if c := d[int(i-deck.Spades)*13+int(j-deck.Ace)]; c.Suit != i || c.Rank != j {
				t.Fatalf("Expected deck to be properly sorted\n Deck:%v\n", d)
			}
		}
	}
}
