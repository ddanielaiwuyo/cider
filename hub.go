package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand/v2"
	"net"
	"zod/protocol"
)

type user struct {
	id   protocol.UserId
	conn net.Conn
}

type PrivateMessage struct {
	Dest protocol.UserId
	Msg  protocol.Response
}

type hub struct {
	register   chan user
	unregister chan protocol.UserId
	broadcast  chan protocol.Response
	private    chan PrivateMessage
}

func (h *hub) run() {
	users := make(map[protocol.UserId]net.Conn)

	for {
		select {
		case user := <-h.register:
			users[user.id] = user.conn
			slog.Info(" [run] new user registered")
		case unregister := <-h.unregister:
			delete(users, unregister)
		case msg := <-h.broadcast:
			content, err := json.Marshal(msg)
			if err != nil {
				slog.Error(" [run] could not marshal res", "", err)
				return
			}
			for id, conn := range users {
				_, err := conn.Write(content)
				if err != nil {
					slog.Error(" [run] could not write to conn", "", err)
					delete(users, id)
					continue
				}
			}

		case privateMsg := <-h.private:
			for id, conn := range users {
				if id == privateMsg.Dest {
					content, err := json.Marshal(privateMsg.Msg)
					if err != nil {
						slog.Error(" [run] could not marshal private msg", "", err)
					}

					if _, err := conn.Write(content); err != nil {
						slog.Error(" [run] could not write to conn", "", err)
					}
				}
			}
		}
	}

}

func newHub() *hub {
	return &hub{
		register:   make(chan user),
		unregister: make(chan protocol.UserId),
		broadcast:  make(chan protocol.Response),
	}
}

func Start() error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}

	hub := newHub()
	go hub.run()

	slog.Info(" server listening on ", "port ", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(" accept err: ", "err", err)
		}

		go HandleConnections(hub, conn)
	}
}

func HandleConnections(hub *hub, conn net.Conn) {
	id := rand.UintN(10)

	hub.register <- user{protocol.UserId(id), conn}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			slog.Error(" [handle_connection] read_err: ", "", err)
			return
		}

		msg := make([]byte, n)
		copy(msg, buffer[:n])

		relayErr := RelayMessage(msg)
		switch {
		case errors.Is(relayErr, ErrUserNotFound):
			res := ServerResponseMsg(ErrUserNotFound.Error(), protocol.ServerErrorResponse)
			if err := writeResponse(conn, res); err != nil {
				slog.Error(" error: ", "", err)
				return
			}

		case errors.Is(relayErr, ErrInternalError):
			res := ServerResponseMsg("Please forgive us bro\n", protocol.ServerErrorResponse)
			if err := writeResponse(conn, res); err != nil {
				slog.Error(" error: ", "", err)
				return
			}

		case errors.Is(relayErr, ErrInvalidMessage):
			res := ServerResponseMsg(ErrInvalidMessage.Error(), protocol.ServerErrorResponse)
			if err := writeResponse(conn, res); err != nil {
				slog.Error(" error: ", "", err)
				return
			}
		default:
			slog.Info(" Yeah, Unexpected error ", "err", err)
		}

	}

}
