package blackjack

import (
	"context"
	"math/rand"
	"time"

	"github.com/mishuk-sk/gopher/deck"
)

type player struct {
	cards []deck.Card
	name  string
	table chan map[string][]deck.Card
	end   <-chan struct{}
	hit   func(string, context.Context) bool
}

func NewPlayer(name string, hit func(string, context.Context) bool) *player {
	return &player{
		name:  name,
		hit:   hit,
		table: make(chan map[string][]deck.Card),
	}
}
func (p *player) OnChange(f func(map[string][]deck.Card)) {
	go func() {
		for {
			select {
			case t := <-p.table:
				f(t)
			case <-p.end:
				return
			}
		}
	}()
}

type dealer struct {
	player
}

func NewDealer() *dealer {
	return &dealer{
		player: player{
			name: "Dealer",
			hit: func(s string, ctx context.Context) bool {
				r := rand.New(rand.NewSource(time.Now().Unix()))
				return r.Intn(2) == 1
			},
		},
	}
}

type table struct {
	players []*player
	deck    []deck.Card
	cards   map[string][]deck.Card
	end     chan struct{}
	dealer  *dealer
}

func NewTable(pl ...*player) *table {
	return &table{
		players: pl,
		deck:    deck.New(deck.Shuffle),
		cards:   nil,
		dealer:  NewDealer(),
	}
}

func (t *table) Start() chan struct{} {
	t.init()
	t.dealer.cards = make([]deck.Card, 0, 2)
	t.dealer.cards = append(t.dealer.cards, t.card())
	t.cards[t.dealer.name] = t.dealer.cards
	t.notify()
	for _, p := range t.players {
		p.cards = make([]deck.Card, 0, 2)
		p.cards = append(p.cards, t.card(), t.card())
		t.cards[p.name] = p.cards
		t.notify()
	}
	for _, p := range t.players {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		for p.hit("Choose hit or fold", ctx) {
			p.cards = append(p.cards, t.card())
			t.cards[p.name] = p.cards
			t.notify()
		}
	}
	for t.dealer.hit("Choose hit or fold", nil) {
		t.dealer.cards = append(t.dealer.cards, t.card())
		t.cards[t.dealer.name] = t.dealer.cards
		t.notify()
	}
	return t.end
}

func (t *table) card() deck.Card {
	d := t.deck
	c := d[0]
	t.deck = d[1:]
	return c
}

func (t *table) init() {
	t.end = make(chan struct{})
	t.cards = make(map[string][]deck.Card, len(t.players))
	for _, p := range t.players {
		t.cards[p.name] = make([]deck.Card, 0, 2)
		p.end = t.end
	}
	t.notify()
}

func (t *table) notify() {
	go func() {
		for _, p := range t.players {
			p.table <- t.cards
		}
	}()
}
