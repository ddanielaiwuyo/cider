package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"zod/logger"
	"zod/protocol"
)

const serverAddr = ":4000"

var ErrInvalidMessage = errors.New("Invalid Message")
var ErrUserNotFound = errors.New("User not found")
var ErrInternalError = errors.New("Internal Server Error")

var connectionPool = make(map[net.Conn]bool)
var connectedUsers = make(map[protocol.UserId]net.Conn)

var jsonLogger = logger.JSONLogger()

type userId uint64

type Request struct {
	Recipient userId `json:"recipient"`
	Msg       string `json:"msg"`
}

type ResponseType int

// type Response struct {
// 	From userId       `json:"from"`
// 	Code ResponseType `json:"responseType"`
// 	Msg  string       `json:"msg"`
// }

func StartServer() error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}

	slog.Info(" server listening on ", "port ", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(" accept err: ", "err", err)
		}

		connectionPool[conn] = true
		fmt.Printf("conn: %s\n", conn.RemoteAddr().String())

		go HandleConnection(conn)
	}
}

func createUser(conn net.Conn) {
	id := len(connectedUsers) + 1
	connectedUsers[protocol.UserId(id)] = conn
}

func removeUser(target net.Conn) {
	for userId, conn := range connectedUsers {
		if conn == target {
			delete(connectedUsers, userId)
			break
		}
	}
}

func broadcast(res protocol.Response) {
	content, err := json.Marshal(res)
	if err != nil {
		slog.Error(" [broadcast] could not marshall response", "err", err)
		return
	}
	for _, conn := range connectedUsers {
		_, err := conn.Write(content)
		if err != nil {
			slog.Error(" [broadcast] could not write to conn", "errr", err)
		}
	}
}

func HandleConnection(conn net.Conn) {
	createUser(conn)

	cu := showConnectedUsers()

	serverRes := ServerResponseMsg(cu, protocol.ServerPaintMessage)
	broadcast(serverRes)

	// writeResponse(conn, res)
	buffer := make([]byte, 1024)
	defer conn.Close()
	defer removeUser(conn)

	for {
		n, err := conn.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			slog.Info(" client disconnected:", "addr:", conn.RemoteAddr().String())
			return
		} else if err != nil && !errors.Is(err, io.EOF) {
			broadcast(serverRes)
			slog.Error(" read_err:", "error", err)
			return
		}

		extractedMessage := buffer[:n]
		go func() {
			err := RelayMessage(extractedMessage)

			switch err != nil {
			case errors.Is(err, ErrUserNotFound):
				res := ServerResponseMsg(ErrUserNotFound.Error(), protocol.ServerErrorResponse)
				if err := writeResponse(conn, res); err != nil {
					slog.Error(" error: ", "", err)
					return
				}

			case errors.Is(err, ErrInternalError):
				res := ServerResponseMsg("Please forgive us bro\n", protocol.ServerErrorResponse)
				if err := writeResponse(conn, res); err != nil {
					slog.Error(" error: ", "", err)
					return
				}

			case errors.Is(err, ErrInvalidMessage):
				res := ServerResponseMsg(ErrInvalidMessage.Error(), protocol.ServerErrorResponse)
				if err := writeResponse(conn, res); err != nil {
					slog.Error(" error: ", "", err)
					return
				}
			default:
				slog.Info(" Yeah, Unexpected error ", "err", err)
			}
		}()
	}
}

func showConnectedUsers() string {
	var msg strings.Builder
	msg.WriteString("    Connected Users,")
	for i := range connectedUsers {
		fmt.Fprintf(&msg, "User %d is Active,", i)
	}

	return msg.String()

}

// The request is parsed into the Request struct format
//
// If the Recipient is not connected at the moment,
// we simply send tell them that we could not find the user
//
// This spec is open to change.
func RelayMessage(request []byte) error {
	req := protocol.Request{}
	if err := json.Unmarshal(request, &req); err != nil {
		slog.Error(" cannot unmarshall:", "err", request)
		return ErrInvalidMessage
	}

	var dest net.Conn

	jsonLogger.Info("incoming", "request", req)

	for id, conn := range connectedUsers {
		if id == req.Recipient {
			dest = conn
			break
		}
	}

	if dest == nil {
		return ErrUserNotFound
	}

	res := protocol.Response{From: req.Recipient, Msg: req.Msg}
	content, err := json.Marshal(res)
	if err != nil {
		slog.Error(" error in marshalling: ", "err", err)
		return ErrInternalError
	}

	if _, err := io.Copy(dest, bytes.NewReader(content)); err != nil {
		return fmt.Errorf(" error in copying to dest: %w", err)
	}
	return nil
}

func writeResponse(conn net.Conn, res protocol.Response) error {
	content, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf(" could not marshall response: %w", err)
	}

	if _, err := conn.Write(content); err != nil {
		return fmt.Errorf(" could not write response: %w", err)
	}

	return nil
}

func main() {
	if err := StartServer(); err != nil {
		slog.Error(" BOMBOCLAT!!\n %s\n", "err", err)
	}
}
