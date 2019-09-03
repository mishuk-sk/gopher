package blackjack

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mishuk-sk/gopher/blackjack/player"
	"github.com/mishuk-sk/gopher/deck"
)

type score map[uuid.UUID][]deck.Card

func (s *score) init(players map[uuid.UUID]*player.Player) {
	m := make(map[uuid.UUID][]deck.Card, len(players))
	for id := range players {
		m[id] = make([]deck.Card, 0, 2)
	}
}

func (s score) add(id uuid.UUID, card deck.Card) (uuid.UUID, error) {
	if _, ok := s[id]; !ok {
		return uuid.Nil, fmt.Errorf("Can't find user with id %s", id)
	}
	s[id] = append(s[id], card)
	sc := calcScore(s[id])
	if sc == 21 {
		return id, nil
	}
	return uuid.Nil, nil
}

func calcScore(cards []deck.Card) int {
	aces := 0
	score := 0
	for _, c := range cards {
		if c.Rank == deck.Ace {
			aces++
		}
		switch c.Rank {
		case deck.Ace, deck.King, deck.Queen, deck.Jake, deck.Ten:
			score += 10
		default:
			score += int(c.Rank)
		}
	}
	for ; score > 21 && aces > 0; aces-- {
		score -= 9
	}
	return score
}

type Table struct {
	ID      uuid.UUID
	Name    string
	Deck    []deck.Card
	Players map[uuid.UUID]*player.Player
}

func NewTable(name string, deckOpts ...deck.Option) *Table {
	return &Table{
		ID:      uuid.New(),
		Name:    name,
		Deck:    deck.New(deckOpts...),
		Players: make(map[uuid.UUID]*player.Player),
	}
}

func (t *Table) Start() (<-chan *player.Player, error) {
	if len(t.Players) == 0 {
		return nil, fmt.Errorf("Can't start blackjack game with 0 players")
	}
	t.score.init(t.Players)

	return nil, nil
}
