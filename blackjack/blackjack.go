package blackjack

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mishuk-sk/gopher/blackjack/player"
	"github.com/mishuk-sk/gopher/deck"
)

type score map[uuid.UUID][]deck.Card

func (s *score) init(ids ...uuid.UUID) {
	m := make(map[uuid.UUID][]deck.Card, len(ids))
	for _, id := range ids {
		m[id] = make([]deck.Card, 0, 2)
	}
	*s = score(m)
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

func (t *Table) Add(players ...*player.Player) {
	for _, p := range players {
		t.Players[p.ID] = p
	}
}

func (t *Table) Start() (<-chan *player.Player, error) {
	if len(t.Players) == 0 {
		return nil, fmt.Errorf("Can't start blackjack game with 0 players")
	}
	t.Deck = deck.Shuffle(t.Deck)
	dealer := player.New("Dealer", func(interface{}, context.Context) {})
	handleGame(t, dealer)
	return nil, nil
}

func handleGame(t *Table, dealer *player.Player) error {

	ids := make([]uuid.UUID, 0, len(t.Players)+1)
	for _, p := range t.Players {
		ids = append(ids, p.ID)
	}
	ids = append(ids, dealer.ID)
	var sc score
	sc.init(ids...)
	card, err := getCard(&t.Deck)
	if err != nil {
		return err
	}
	_, err = sc.add(dealer.ID, card)
	if err != nil {
		panic(err)
	}
	notify(sc, t.Players)
	for i := 0; i < 2; i++ {
		for _, p := range t.Players {
			card, err := getCard(&t.Deck)
			if err != nil {
				panic(err)
			}
			_, err = sc.add(p.ID, card)
			if err != nil {
				panic(err)
			}
			notify(sc, t.Players)
		}
	}
	return nil
}

func notify(data interface{}, players map[uuid.UUID]*player.Player) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	for _, p := range players {
		p.Notify(data, ctx)
	}
}

func getCard(d *[]deck.Card) (deck.Card, error) {
	nd := *d
	if len(nd) == 0 {
		return deck.Card{}, fmt.Errorf("Can't get card from empty deck")
	}
	card := nd[0]
	if len(nd) > 1 {
		nd = nd[1:]
	} else {
		nd = []deck.Card{}
	}
	*d = nd
	return card, nil
}
