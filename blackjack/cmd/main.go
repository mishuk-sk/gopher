package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mishuk-sk/gopher/blackjack"
	"github.com/mishuk-sk/gopher/blackjack/player"
)

func main() {
	players := make([]*player.Player, 1)
	for i := range players {
		players[i] = player.New(fmt.Sprintf("Player %d", i), func(data interface{}, ctx context.Context) {
			fmt.Println(data)
		})
	}
	table := blackjack.NewTable("Table")
	table.Add(players...)
	table.Start()
	<-time.After(time.Second * 3)
}
