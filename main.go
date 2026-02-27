package main

import (
	"context"
	"log"
	"github.com/persona-mp3/internal/server"
)

func main() {
	manager := server.NewManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go manager.Listen(ctx)

	if err := server.Start(manager); err != nil {
		log.Fatal(err)
	}
}
