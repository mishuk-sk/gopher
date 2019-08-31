package deck

import "fmt"

//go:generate stringer -type=Suit

//Suit represents cards suit
type Suit int

//Constants block represents cards suits
const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
)

//go:generate stringer -type=Rank

//Rank represents type for card rank
type Rank int

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

//Card represents card by 2 values: Suit and Rank.
//Suit (0-3) represents card Suit
//Rank (0-12) represents card value.
//Both are defined as constants in deck package
type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

//New creates new deck of 52 cards sorted in default order.
func New() []Card {
	deck := make([]Card, 0, 52)
	for i := Spades; i <= Hearts; i++ {
		for j := Ace; j <= King; j++ {
			deck = append(deck, Card{
				Suit: i,
				Rank: j,
			})
		}
	}
	return deck
}
