package blackjack

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mishuk-sk/gopher/deck"
)

type Player struct {
	Cards        []deck.Card
	Notification func(msg interface{}, ctx context.Context)
	Hit          func() bool
}

func NewPlayer(notification func(msg interface{}, ctx context.Context), hit func() bool) *Player {
	return &Player{
		Notification: notification,
		Hit:          hit,
	}
}

//GiveCard adds card to player's cards and returns current score
func GiveCard(p *Player, card deck.Card) int {
	p.Cards = append(p.Cards, card)
	return calcScore(p.Cards)
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

func Notify(p *Player, msg interface{}, ctx context.Context) {
	//Double goroutine to handle p.Notification cancel correct, when not handled inside
	//FIXME probably leaking goroutine
	go func() {
		p.Notification(msg, ctx)
	}()
}

type Game struct {
	Deck    deck.Deck
	Players []*Player
}

func NewGame(name string, deckOpts ...deck.Option) *Game {
	return &Game{
		Deck:    deck.New(deckOpts...),
		Players: []*Player{},
	}
}

func (t *Game) Start() (<-chan *Player, error) {
	if len(t.Players) == 0 {
		return nil, fmt.Errorf("Can't start blackjack game with 0 players")
	}
	t.Deck = deck.Shuffle(t.Deck)
	dealer := NewPlayer(func(interface{}, context.Context) {}, func() bool {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		return r.Intn(3) == 0
	})
	handleGame(t, dealer)
	return nil, nil
}

func handleGame(g *Game, dealer *Player) error {
	deck := &g.Deck
	card, err := deck.Card()
	if err != nil {
		return err
	}
	GiveCard(dealer, card)
	for i := 0; i < 2; i++ {
		for _, p := range g.Players {
			card, err := deck.Card()
			if err != nil {
				return err
			}
			GiveCard(p, card)
		}
	}
	notifyPlayers("Some info", g.Players...)
}

func notifyPlayers(data interface{}, players ...*Player) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	for _, p := range players {
		go func() {
			p.Notification(data, ctx)
		}()
	}
}
