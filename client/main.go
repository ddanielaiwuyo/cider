package main

import (
	"context"
	"fmt"
	"log"
	"net"
)

func runner() {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
		return
	}

	stdin := make(chan string)
	serverCh := make(chan []byte)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go readFromStdin(ctx, stdin)

	go func() {
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
				return
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
	fmt.Printf("message from stdin: %s\n", msg)

	if _, err := conn.Write([]byte(msg)); err != nil {
		return fmt.Errorf(" write error: %w", err)
	}

	return nil
}

// this is where part of the game logic will reside
func handleServerResponse(response []byte) {
	fmt.Printf(" # : %s\n", response)
}

func main(){
	runner()
}
