package blackjack

import (
	"context"
	"fmt"
	"time"

	"github.com/mishuk-sk/gopher/deck"
)

type Player struct {
	Cards []deck.Card
}

func NewPlayer() Player {
	return Player{
		Cards: []deck.Card{},
	}
}

//GiveCard adds card to player's cards and returns current score
func (p *Player) GiveCard(card deck.Card) int {
	p.Cards = append(p.Cards, card)
	return calcScore(p.Cards)
}
func (p Player) Hit() bool {
	fmt.Println("Choose to hit or stay")
	var in string
	fmt.Scanf("%s", &in)
	switch in {
	case "h":
		return true
	case "s":
		return false
	}
	return false
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
func (p Player) Notify(data interface{}) {
	fmt.Println(data)
}

const (
	StatePlayerTurn = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck      deck.Deck
	State     int
	Players   []Player
	Dealer    Player
	curPlayer int
}

func (gs *GameState) CurrentPlayer() *Player {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Players[gs.curPlayer]
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("Not a turn")
	}
}

func Start() {
	gs := GameState{
		Deck:    deck.New(deck.AddDecks(2), deck.Shuffle),
		Players: []Player{NewPlayer(), NewPlayer(), NewPlayer()},
		Dealer:  NewPlayer(),
	}
	gs = Deal(gs)
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	for gs.State == StatePlayerTurn {
		for hitStay(ctx, *gs.CurrentPlayer()) {
			gs = Hit(gs)
			close()
			ctx, close = context.WithTimeout(context.Background(), time.Second*5)
		}
		gs = Stay(gs)
	}
	close()
	for gs.State == StateDealerTurn {
		gs.State = StateHandOver
	}
	fmt.Println(gs.Players)
	fmt.Println(gs.Dealer)
}

func Deal(gs GameState) GameState {
	dealQueue := make([]*Player, 0, len(gs.Players)+1)
	dealQueue = append(dealQueue, &gs.Dealer)
	for i := range gs.Players {
		dealQueue = append(dealQueue, &gs.Players[i])
	}
	for i := 0; i < 2; i++ {
		for _, p := range dealQueue {
			c, _ := gs.Deck.Card()
			p.GiveCard(c)
		}
	}
	gs.State = StatePlayerTurn
	return gs
}
func hitStay(ctx context.Context, p Player) bool {
	choice := make(chan bool, 1)
	go func() {
		choice <- p.Hit()
	}()
	select {
	case c := <-choice:
		return c
	case <-ctx.Done():
		return false
	}
}

func Hit(gs GameState) GameState {
	c, _ := gs.Deck.Card()
	gs.CurrentPlayer().GiveCard(c)
	return gs
}

func Stay(gs GameState) GameState {
	if gs.curPlayer == len(gs.Players)-1 {
		gs.State = StateDealerTurn
	} else {
		gs.curPlayer++
	}
	return gs
}
