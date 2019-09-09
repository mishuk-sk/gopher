package blackjack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mishuk-sk/gopher/deck"
)

type Participant interface {
	Hit() bool
	Notify(interface{})
}
type Hand []deck.Card

func NewHand() Hand {
	return []deck.Card{}
}
func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range strs {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

type Player struct {
	Hand
	Name string
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
func (p Player) Notify(data interface{}) {
	fmt.Println(data)
}
func calcScore(cards []deck.Card) int {
	aces := 0
	score := 0
	for _, c := range cards {
		if c.Rank == deck.Ace {
			aces++
			score++
		}
		switch c.Rank {
		case deck.Ace, deck.King, deck.Queen, deck.Jake, deck.Ten:
			score += 10
		default:
			score += int(c.Rank)
		}
	}
	for ; score > 21 && aces > 0; aces-- {
		score -= 10
	}
	return score
}

type Dealer struct {
	Hand
}

func (d Dealer) Hit() bool {
	score := calcScore(d.Hand)
	return score < 16 || (score == 17 && hasAce(d.Hand))
}
func (d Dealer) Notify(data interface{}) {
	fmt.Println(data)
}

func hasAce(cards []deck.Card) bool {
	for i := range cards {
		if cards[i].Rank == deck.Ace {
			return true
		}
	}
	return false
}

//GiveCard adds card to player's cards and returns current score
func GiveCard(h *Hand, card deck.Card) {
	*h = append(*h, card)
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
	Dealer    Dealer
	curPlayer int
}

func (gs *GameState) CurrentPlayer() Participant {
	switch gs.State {
	case StatePlayerTurn:
		return gs.Players[gs.curPlayer]
	case StateDealerTurn:
		return gs.Dealer
	default:
		panic("Not a turn")
	}
}
func (gs *GameState) CurrentHand() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Players[gs.curPlayer].Hand
	case StateDealerTurn:
		return &gs.Dealer.Hand
	default:
		panic("Not a turn")
	}
}

func Start() {
	participants := []Player{
		Player{
			Hand: NewHand(),
			Name: "Player1",
		},
		Player{
			Hand: NewHand(),
			Name: "Player2",
		},
	}
	gs := GameState{
		Deck:    deck.New(deck.AddDecks(2), deck.Shuffle),
		Players: participants,
		Dealer:  Dealer{NewHand()},
	}
	gs = Deal(gs)

	for gs.State == StatePlayerTurn {
		ctx, close := context.WithTimeout(context.Background(), time.Second*5)

		p := gs.CurrentPlayer().(Player)
		p.Notify(fmt.Sprintf("%s: %s\n Score:%d\n", p.Name, p.Hand, calcScore(p.Hand)))
		for hitStay(ctx, gs.CurrentPlayer()) {
			gs = Hit(gs)
			p := gs.CurrentPlayer().(Player)
			p.Notify(fmt.Sprintf("%s: %s\n Score:%d\n", p.Name, p.Hand, calcScore(p.Hand)))
			close()
			ctx, close = context.WithTimeout(context.Background(), time.Second*5)
		}
		gs = Stay(gs)
		close()
	}
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	for gs.State == StateDealerTurn {

		for hitStay(ctx, gs.CurrentPlayer()) {
			gs = Hit(gs)
			p := gs.CurrentPlayer().(Dealer)
			p.Notify(fmt.Sprintf("%s: %s\n Score:%d\n", "Dealer", p.Hand, calcScore(p.Hand)))
			close()
			ctx, close = context.WithTimeout(context.Background(), time.Second*5)
		}

		p := gs.CurrentPlayer().(Dealer)
		gs = Stay(gs)
		p.Notify(fmt.Sprintf("%s: %s\n Score:%d\n", "Dealer", p.Hand, calcScore(p.Hand)))
	}
	close()
	gs = EndGame(gs)
}

func EndGame(gs GameState) GameState {
	type PlayerScore struct {
		name  string
		score int
	}
	scores := make([]PlayerScore, len(gs.Players)+1)
	i := 0
	max := 0
	maxI := []int{}
	for i = 0; i < len(gs.Players); i++ {
		scores[i] = PlayerScore{
			gs.Players[i].Name,
			calcScore(gs.Players[i].Hand),
		}
		sc := truncateScore(scores[i].score)
		if sc >= max {
			maxI = append(maxI, i)
			max = sc
		}
	}
	if ds := calcScore(gs.Dealer.Hand); max < ds {
		fmt.Printf("Dealer wins with score %d.\n", ds)
	} else if ds == max {
		fmt.Printf("Dealers has score %d, equivalent to players:\n", ds)
		for _, i := range maxI {
			fmt.Printf("\t%s\n", gs.Players[i].Name)
		}
	} else {
		if len(maxI) > 1 {
			fmt.Printf("Players same score have score %d:\n", max)
			for _, i := range maxI {
				fmt.Printf("\t%s\n", gs.Players[i].Name)
			}
		} else {
			fmt.Printf("Player %s wins with score %d:\n", gs.Players[maxI[0]].Name, max)
		}
	}
	return gs
}

func truncateScore(s int) int {
	if s > 21 {
		return 0
	}
	return s
}

func Deal(gs GameState) GameState {
	dealQueue := make([]*Hand, 0, len(gs.Players)+1)
	dealQueue = append(dealQueue, &(gs.Dealer.Hand))
	for i := range gs.Players {
		dealQueue = append(dealQueue, &gs.Players[i].Hand)
	}
	for i := 0; i < 2; i++ {
		for _, h := range dealQueue {
			c, _ := gs.Deck.Card()
			GiveCard(h, c)
		}
	}
	gs.State = StatePlayerTurn
	return gs
}

func hitStay(ctx context.Context, p Participant) bool {
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
	h := gs.CurrentHand()
	GiveCard(h, c)
	return gs
}

func Stay(gs GameState) GameState {
	if gs.State == StateDealerTurn {
		gs.State = StateHandOver
	} else if gs.curPlayer == len(gs.Players)-1 {
		gs.State = StateDealerTurn
	} else {
		gs.curPlayer++
	}
	return gs
}
