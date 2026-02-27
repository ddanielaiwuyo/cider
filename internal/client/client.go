package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/persona-mp3/shared"
)

const serverAddr = "localhost:4000"

func connect() error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("could not connect server: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stdin := readFromStdin(ctx)
	server := readFromServer(ctx, conn)
	for {
		select {
		case msg, ok := <-stdin:
			if !ok {
				return fmt.Errorf("stdin has been closed!")
			}

			// fmt.Println(" *stdin: ", msg)
			if err := writeToServer(msg, conn); err != nil {
				return err
			}

		case res, ok := <-server:
			if !ok {
				return fmt.Errorf("server channel has been closed!")
			}

			parseServerResponse(res)
		}
	}
}

func parseServerResponse(res shared.Message) {
	fmt.Printf("  #%d:  %2s\n", res.From, res.Content)
}

func main() {
	if err := connect(); err != nil {
		log.Fatal(err)
	}
}
