package player

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Player struct {
	ID           uuid.UUID
	Name         string
	Notification func(msg interface{}, ctx context.Context)
}

func New(name string, notification func(msg interface{}, ctx context.Context)) *Player {
	return &Player{
		ID:           uuid.New(),
		Name:         name,
		Notification: notification,
	}
}

func (p *Player) Notify(msg interface{}, ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	//Double goroutine to handle p.Notification cancel correct, when not handled inside
	//FIXME probably leaking goroutine
	fmt.Println(p.Name)
	p.Notification(msg, ctx)
	close(done)
	return done
}
