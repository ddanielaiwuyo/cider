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
var connectionPool = make(map[net.Conn]bool)

const serverAddr = ":4000"

func StartServer() error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}

	log.Println(" server listening on port ", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err: %s\n", err)
		}

		connectionPool[conn] = true
		fmt.Printf("conn: %s\n", conn.RemoteAddr().String())

		go HandleConnection(conn)
	}
}



func HandleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	defer conn.Close()
	defer delete(connectionPool, conn)

	peer := CreatePeer(conn)
	for {
		n, err := peer.conn.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			slog.Info(" client disconnected:", "addr:", peer.addr)
			return
		} else if err != nil && !errors.Is(err, io.EOF) {
			slog.Error(" read_err:", "error", err)
			return
		}

		extractedMessage := buffer[:n]
		msgType, content := protocol.ParseMessage(extractedMessage)
		switch msgType {
		case protocol.ConnectTo:
			peer.ConnectTo(content)
		// default:
		// 	go Broadcast(extractedMessage)
		}

		// if unknown message type, either echo back to them
		// or just brodcrast it and tell other people

	}
}

func Broadcast(msg []byte) {
	fmt.Println("broadcasting")
	for conn := range connectionPool {
		_, err := conn.Write(msg)
		if err != nil {
			log.Printf(" write_err %s: %s\n", conn.RemoteAddr().String(), err)
			delete(connectionPool, conn)
		}
	}
}

func main() {
	if err := StartServer(); err != nil {
		log.Fatalf(" BOMBOCLAT!!\n %s\n", err)
	}
}
