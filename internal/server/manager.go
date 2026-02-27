package server

import (
	"context"
	"log/slog"
)

type manager struct {
	register chan client
	remove   chan int
	deliver  chan Message
	find     chan int
}

func NewManager() *manager {
	return &manager{
		register: make(chan client),
		remove:   make(chan int),
		deliver:  make(chan Message),
	}
}

func (m *manager) Listen(ctx context.Context) {
	for {
		select {
		case client := <-m.register:
			slog.Info("adding new client", "", client.id)
			connectedUsers[client.id] = client.conn
		case id := <-m.remove:
			slog.Info("removing client with id:", "", id)
			delete(connectedUsers, id)
		case msg := <-m.deliver:
			sendMessage(m, msg)
		case <-ctx.Done():
			slog.Info("exiting manager:", "", ctx.Err().Error())
			return
		}
	}
}
