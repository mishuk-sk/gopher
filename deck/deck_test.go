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
	for i := deck.Spade; i <= deck.Heart; i++ {
		for j := deck.Ace; j <= deck.King; j++ {
			if c := d[int(i-deck.Spade)*13+int(j-deck.Ace)]; c.Suit != i || c.Rank != j {
				t.Fatalf("Expected deck to be properly sorted\n Deck:%v\n", d)
			}
		}
	}
}

func TestSort(t *testing.T) {
	d := deck.New()
	deck.Sort(d, func(i, j int) bool {
		return d[i].Rank < d[j].Rank
	})
	for i := deck.Ace; i <= deck.King; i++ {
		for j := 0; j < 4; j++ {
			if index := int(i-deck.Ace)*4 + j; d[index].Rank != i {
				t.Fatalf("Expected to get %s on position %d, but got %s\n", i, index, d[index].Rank)
			}
		}
	}
}

func TestDefaultSort(t *testing.T) {
	d := deck.New()
	deck.DefaultSort(d)
	for i := deck.Spade; i <= deck.Heart; i++ {
		for j := deck.Ace; j <= deck.King; j++ {
			if c := d[int(i-deck.Spade)*13+int(j-deck.Ace)]; c.Suit != i || c.Rank != j {
				t.Fatalf("Expected deck to be properly sorted\n Deck:%v\n", d)
			}
		}
	}
}
