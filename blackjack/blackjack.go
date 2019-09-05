package blackjack

import (
	"fmt"

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

/*
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
*/
const (
	StatePlayerTurn = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   deck.Deck
	State  int
	Player Player
	Dealer Player
}

func (gs *GameState) CurrentPlayer() *Player {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("Not a turn")
	}
}

func Start() {
	gs := GameState{
		Deck:   deck.New(deck.AddDecks(2), deck.Shuffle),
		Player: NewPlayer(),
		Dealer: NewPlayer(),
	}
	gs = Deal(gs)
	gs.Player.Notify(gs.Dealer)
	for gs.State == StatePlayerTurn {
		fmt.Println("Choose to hit or stay")
		var in string
		fmt.Scanf("%s", &in)
		switch in {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stay(gs)
		}
	}
	for gs.State == StateDealerTurn {
		gs = Stay(gs)
	}
	fmt.Println(gs.Player)
	fmt.Println(gs.Dealer)
}

func Deal(gs GameState) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range []*Player{&gs.Player, &gs.Dealer} {
			c, _ := gs.Deck.Card()
			p.GiveCard(c)
		}
	}
	gs.State = StatePlayerTurn
	return gs
}

func Hit(gs GameState) GameState {
	c, _ := gs.Deck.Card()
	gs.CurrentPlayer().GiveCard(c)
	return gs
}

func Stay(gs GameState) GameState {
	gs.State++
	return gs
}
