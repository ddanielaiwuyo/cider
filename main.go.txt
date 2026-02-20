package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"zod/protocol"
)

const socketAddr = ":4000"

var connections = make(map[net.Conn]bool)

func StartServer() error {
	listener, err := net.Listen("tcp", socketAddr)
	if err != nil {
		return err
	}

	log.Println(" server listening on port ", socketAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err: %s\n", err)
		}

		connections[conn] = true
		fmt.Printf("conn: %s\n", conn.RemoteAddr().String())

		go handle_connections(conn)
	}
}

func handle_connections(conn net.Conn) {
	buffer := make([]byte, 1024)
	defer conn.Close()
	defer delete(connections, conn)

	client := CreateClient(conn)
	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info(" client-disconnected, eof")
			} else {
				log.Println(" read_err: ", err)
			}
			return
		}

		fmt.Printf(" %s\n", buffer[:n])
		msgType, content := protocol.ParseMessage(buffer[:n])
		switch msgType {
		case protocol.ConnectTo:
			fmt.Printf(" sender: %s requested to connect to: %s\n", client.conn.RemoteAddr().String(), content)
			err := handleTwoWayConnect(client, content)
			if err != nil {
				slog.Error("HandlingTwoWayConnect:", "error", err)
				return
			}
		}

		errCh := make(chan error)
		go func() {
			errCh <- Broadcast(client, buffer[:n])
		}()

		errs := <-errCh
		if errs != nil {
			log.Println(err)
		}

	}
}

func Broadcast(client Client, msg []byte) error {
	for conn := range connections {
		if conn == client.conn {
			continue
		}
		_, err := conn.Write(msg)
		if err != nil {
			return fmt.Errorf(" write_err %s: %w", conn.RemoteAddr().String(), err)
		}
	}
	return nil
}

// so we need to look for a way to
// allow the user to connect to someone
// in particular

func main() {
	fmt.Printf("\n Welcome back apple cider vinegar\n")
	if err := StartServer(); err != nil {
		log.Fatal(err)
	}

}

func handleTwoWayConnect(client Client, destAddr string) error {
	_, err := net.ResolveTCPAddr("tcp", destAddr)
	if err != nil {
		return fmt.Errorf("could not parse destAddr %w", err)
	}

	if err := client.ConnectTo(destAddr); err != nil {
		return err
	}
	return  nil
}

