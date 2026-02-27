package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/google/uuid"
)

const serverAddr = ":4000"

var ErrMalformedMessage = errors.New("Message is invalid")
var ErrContactUser = errors.New("Could not contact user")
var ErrInternalError = errors.New("Please forgive us bro")

var NotFoundResponse = Message{
	From:    0,
	Content: ErrContactUser.Error(),
}

var connectedUsers = make(map[int]net.Conn)

type client struct {
	id   int
	conn net.Conn
}

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
	start := 1000
	for {
		select {
		case client := <-m.register:
			slog.Info("adding new client", "", start)
			connectedUsers[start] = client.conn
			start += 100
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

func Start(mgr *manager) error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("could not start server %w", err)
	}

	log.Printf("server active on localhost%s\n", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("accept_connection err: %w ", "", err)
			continue
		}

		go HandleConnection(mgr, conn)
	}
}

type Message struct {
	From    int    `json:"from"`
	Content string `json:"content"`
	Dest    int    `json:"dest"`
}

func sendMessage(mgr *manager, msg Message) {
	var notFoundRes = NotFoundResponse
	notFoundRes.Dest = msg.From

	id := msg.Dest
	destConn, destfound := connectedUsers[id]
	senderConn, found := connectedUsers[msg.From]
	if !found {
		slog.Info("sender not recognised", "id", msg.From, "found", found)
		return
	}

	if !destfound {
		slog.Error("dst conn not found", "", id)

		res, err := toJson(notFoundRes)
		if err != nil {
			log.Println(err)
			return
		}

		io.Copy(senderConn, bytes.NewReader(res))
		return
	}
	fmt.Println("destConn found", destConn)

	content, err := toJson(msg)
	if err != nil {
		slog.Error("", "", err)
		return
	}
	if _, err := io.Copy(destConn, bytes.NewReader(content)); err != nil {
		slog.Error("error writing to dest-conn", "", err)
		mgr.remove <- id
		return
	}

	log.Println(" [success] message sent")
}

func HandleConnection(mgr *manager, conn net.Conn) {
	var newClient = client{id: int(uuid.New().ID()), conn: conn}
	var welcomeResponse = Message{From: 0, Content: "Welcome to CiderVine", Dest: newClient.id}

	mgr.register <- newClient
	defer func() {
		mgr.remove <- newClient.id
	}()

	buff := make([]byte, 1024)
	content, err := toJson(welcomeResponse)
	if err != nil {
		slog.Error("", "", err)
		return
	}

	conn.Write(content)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			slog.Error("read error from client:", "", err)
			return
		}

		dest := make([]byte, n)
		copy(dest, buff[:n])

		request := Message{}
		if err := json.Unmarshal(dest, &request); err != nil {
			slog.Error("could not parse request, malformed", "err", err)
			mgr.deliver <- Message{Dest: newClient.id, Content: ErrMalformedMessage.Error(), From: 0}
			continue
		}

		mgr.deliver <- request
	}
}
