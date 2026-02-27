package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
)

const serverAddr = ":4000"
const serverId = 0

var ErrMalformedMessage = errors.New("Message is invalid")
var ErrContactUser = errors.New("Could not contact user")
var ErrInternalError = errors.New("Please forgive us bro")

var NotFoundResponse = Message{
	From:    serverId,
	Content: ErrContactUser.Error(),
}

var connectedUsers = make(map[int]net.Conn)

type client struct {
	id   int
	conn net.Conn
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

		go handleConnection(mgr, conn)
	}
}

type MessageType int

const (
	PaintMessage MessageType = iota
	ChatMessage
	GameMessage
)

type Message struct {
	From        int         `json:"from"`
	MessageType MessageType `json:"messageType"`
	Content     string      `json:"content"`
	Dest        int         `json:"dest"`
}

func handleConnection(mgr *manager, conn net.Conn) {
	var paintMsg = createPaintMessage()
	var newClient = client{
		id:   len(connectedUsers) + 1,
		conn: conn,
	}

	var welcomeResponse = Message{
		From:        serverId,
		MessageType: PaintMessage,
		Content:     fmt.Sprintf("Welcome to CiderVine;%s;YourId:%d;", paintMsg, newClient.id),
		Dest:        newClient.id,
	}

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
			// mgr.deliver <- Message{
			// 	Dest:    newClient.id,
			// 	Content: ErrMalformedMessage.Error(),
			// 	From:    serverId,
			// }
			mgr.deliver <- createMalformedMessage(newClient.id)
			continue
		}

		mgr.deliver <- request
	}
}

func sendMessage(mgr *manager, msg Message) {
	var notFoundRes = NotFoundResponse
	notFoundRes.Dest = msg.From

	id := msg.Dest
	senderId := msg.From

	destConn, destfound := connectedUsers[id]
	senderConn, found := connectedUsers[senderId]

	if !found {
		slog.Info("sender not recognised", "id", serverId, "found", found)
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
