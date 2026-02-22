package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"zod/protocol"
)

func runner() {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
		return
	}

	stdin := make(chan string)
	serverCh := make(chan []byte)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go readFromStdin(ctx, stdin)

	go func() {
		defer conn.Close()
		defer close(serverCh)
		if err := readFromServer(ctx, conn, serverCh); err != nil {
			log.Println(err)
			return
		}
	}()

	for {
		select {
		case val, ok := <-stdin:
			if !ok {
				log.Println(" stdin channel closed")
				return
			}
			// temporary hacking
			if val == "q" {
				fmt.Println(" user quitting")
				close(stdin)
				break
			}

			handleStdinMessage(conn, val)
		case response, ok := <-serverCh:
			if !ok {
				log.Println(" server channel closed")
				return
			}
			handleServerResponse(response)
		}
	}

}

func handleStdinMessage(conn net.Conn, msg string) error {
	if _, err := conn.Write([]byte(msg)); err != nil {
		return fmt.Errorf(" write error: %w", err)
	}

	return nil
}

// this is where part of the game logic will reside
func handleServerResponse(response []byte) {
	var res protocol.Response
	if err := json.Unmarshal(response, &res); err != nil {
		log.Println(" could  not parse server's response")
		return
	}

	clearTerminal()
	msg := res.Msg
	sender := res.From

	if int(sender) == 4000 && res.Code == protocol.ServerPaintMessage {
		displayConnectedUsers(msg)
		return
	}

	fmt.Printf("    \t#%d: %4s\n", sender, msg)

}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func displayConnectedUsers(msg string) {
	for msg := range strings.SplitSeq(msg, ",") {
		fmt.Printf("    \t%s\n", msg)
	}
}
func main() {
	runner()
}
