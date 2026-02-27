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

var connectedUsers = make(map[int]net.Conn)

type client struct {
	id   int
	conn net.Conn
}

type manager struct {
	register chan client
	remove   chan int
	deliver  chan Message
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
			slog.Info("adding new client", "", client)
			// connectedUsers.Store(client.id, client.conn)
			connectedUsers[client.id] = client.conn
		case id := <-m.remove:
			slog.Info("removing client with id:", "", id)
			// connectedUsers.Delete(id)
			delete(connectedUsers, id)
		case msg := <-m.deliver:
			go handleMessage(m, msg)
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

func handleMessage(mgr *manager, msg Message) {
	destId := msg.Dest
	destConn, found := connectedUsers[destId]
	if !found {
		slog.Info("could not send to user", "", destId)
		return
	}

	content, err := json.Marshal(msg)
	if err != nil {
		slog.Error("could not marshall msg", "", err)
		return
	}

	if destConn == nil {
		return
	}
	if _, err := io.Copy(destConn, bytes.NewReader(content)); err != nil {
		slog.Error("write error to dest", "", err)
		mgr.remove <- destId
		return
	}

	slog.Info("successfully written to dest")
}

func HandleConnection(mgr *manager, conn net.Conn) {
	newClient := client{id: int(uuid.New().ID()), conn: conn}
	mgr.register <- newClient
	defer func() {
		mgr.remove <- newClient.id
	}()

	buff := make([]byte, 1024)
	conn.Write([]byte("hello world"))
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
