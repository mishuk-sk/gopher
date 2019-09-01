package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mishuk-sk/gopher/blackjack"
	"github.com/mishuk-sk/gopher/deck"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	player := blackjack.NewPlayer("Player1", func(str string, ctx context.Context) bool {

		n := r.Intn(3)
		return n == 1
	})
	player.OnChange(func(m map[string][]deck.Card) {
		for k, v := range m {
			fmt.Printf("%s: %s\n", k, v)
		}
	})
	table := blackjack.NewTable(player)
	table.Start()
	<-time.After(time.Second * 3)
}
