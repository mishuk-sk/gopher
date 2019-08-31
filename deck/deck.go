package deck

import (
	"fmt"
	"sort"
)

//go:generate stringer -type=Suit

//Suit represents cards suit
type Suit uint8

//Constants block represents cards suits
const (
	Spade Suit = iota
	Diamond
	Club
	Heart
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

//go:generate stringer -type=Rank

//Rank represents type for card rank
type Rank uint8

//Constants block represents cards ranks
const (
	Joker Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jake
	Queen
	King
)
const (
	minRank = Ace
	maxRank = King
)

//Card represents card by 2 values: Suit and Rank.
//Suit (0-3) represents card Suit
//Rank (0-12) represents card value.
//Both are defined as constants in deck package
type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	if c.Rank == Joker {
		return "Joker"
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}

//New creates new deck of 52 cards sorted in default order.
func New() []Card {
	deck := make([]Card, 0, 52)
	for _, i := range suits {
		for j := minRank; j <= maxRank; j++ {
			deck = append(deck, Card{
				Suit: i,
				Rank: j,
			})
		}
	}
	return deck
}

//Sort sorts deck in place using provided less(i, j int)bool function
func Sort(deck []Card, less func(i, j int) bool) {
	sort.Slice(deck, less)
}

//DefaultSort sorts deck in default order
func DefaultSort(deck []Card) {
	Sort(deck, func(i, j int) bool {
		return absRank(deck[i]) < absRank(deck[j])
	})
}

//Filter returns function to filter cards according to provided keep function
func Filter(keep func(Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var deck []Card
		for _, c := range cards {
			if keep(c) {
				deck = append(deck, c)
			}
		}
		return deck
	}
}

//AddJokers returns function to add n jokers to the end of deck
func AddJokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		deck := make([]Card, len(cards), len(cards)+n)
		copy(deck, cards)
		for i := 0; i < n; i++ {
			deck = append(deck, Card{
				Rank: Joker,
				Suit: Spade,
			})
		}
		return deck
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}
