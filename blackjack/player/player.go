package player

import (
	"context"

	"github.com/mishuk-sk/gopher/deck"
)

type Player struct {
	Cards        []deck.Card
	Notification func(msg interface{}, ctx context.Context)
	Hit          func() bool
}

func New(notification func(msg interface{}, ctx context.Context), hit func() bool) *Player {
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
