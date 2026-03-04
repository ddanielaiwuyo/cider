package main

import (
	"context"
	"github.com/persona-mp3/internal/server"
	"log"
)

func main() {
	manager := server.NewManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go manager.Listen(ctx)

	if err := server.RunServer(manager); err != nil {
		log.Fatal(err)
	}
}
