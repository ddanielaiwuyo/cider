package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
)

const socketAddr = ":4000"

var connections = make(map[net.Conn]bool)

func start_server() error {
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

		go handle_connections(conn)
	}
}

func broadcast(sender net.Conn, msg []byte) error {
	for conn := range connections {
		if conn == sender {
			continue
		}
		_, err := conn.Write(msg)
		if err != nil {
			return fmt.Errorf(" write_err %s: %w", conn.RemoteAddr().String(), err)
		}
	}
	return nil
}

func handle_connections(conn net.Conn) {
	buffer := make([]byte, 1024)
	defer conn.Close()
	defer delete(connections, conn)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info(" client-disconnected, eof")
			} else {
				log.Println(" read_err: ", err)
			}
			return
		}

		fmt.Printf(" %s\n", buffer[:n])

		errCh := make(chan error)
		go func() {
			errCh <- broadcast(conn, buffer[:n])
		}()

		errs := <-errCh
		if errs != nil {
			log.Println(err)
		}

	}
}

func main() {
	fmt.Printf("\n Welcome back apple cider vinegar\n")
	if err := start_server(); err != nil {
		log.Fatal(err)
	}

}
