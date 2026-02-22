package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"zod/logger"
)

const serverAddr = ":4000"

var ErrInvalidMessage = errors.New("Invalid Message")
var ErrUserNotFound = errors.New("User not found")
var ErrInternalError = errors.New("Internal Server Error")

var connectionPool = make(map[net.Conn]bool)
var connectedUsers = make(map[userId]net.Conn)

var jsonLogger = logger.JSONLogger()

type userId uint64

type Request struct {
	UserId userId `json:"userId"`
	Msg    string `json:"msg"`
}

type Response struct {
	From userId `json:"from"`
	Msg  string `json:"msg"`
}

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
	connectedUsers[userId(id)] = conn
}

func removeUser(target net.Conn) {
	for userId, conn := range connectedUsers {
		if conn == target {
			delete(connectedUsers, userId)
			break
		}
	}
}

func HandleConnection(conn net.Conn) {
	createUser(conn)
	buffer := make([]byte, 1024)
	defer conn.Close()
	defer removeUser(conn)

	for {
		n, err := conn.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			slog.Info(" client disconnected:", "addr:", conn.RemoteAddr().String())
			return
		} else if err != nil && !errors.Is(err, io.EOF) {
			slog.Error(" read_err:", "error", err)
			return
		}

		extractedMessage := buffer[:n]
		go func() {
			err := RelayMessage(extractedMessage)

			switch err != nil {
			case errors.Is(err, ErrUserNotFound):
				if err := writeResponseToConn(conn, CreateResponse(ErrUserNotFound.Error())); err != nil {
					slog.Error(" error: ", "", err)
					return
				}
			case errors.Is(err, ErrInternalError):
				if err := writeResponseToConn(conn, CreateResponse("Please forgive us bro\n")); err != nil {
					slog.Error(" error: ", "", err)
					return
				}

			case errors.Is(err, ErrInvalidMessage):
				if err := writeResponseToConn(conn, CreateResponse(ErrInvalidMessage.Error())); err != nil {
					slog.Error(" error: ", "", err)
					return
				}
			default:
				slog.Info(" Yeah, Unexpected error ", "err", err)
			}
		}()
	}
}

func RelayMessage(request []byte) error {
	req := Request{}
	if err := json.Unmarshal(request, &req); err != nil {
		slog.Error(" cannot unmarshall:", "err", request)
		return ErrInvalidMessage
	}

	var dest net.Conn

	jsonLogger.Info("incoming", "request", req)

	for id, conn := range connectedUsers {
		if id == req.UserId {
			dest = conn
			break
		}
	}

	if dest == nil {
		return ErrUserNotFound
	}

	res := Response{From: req.UserId, Msg: req.Msg}
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

func writeResponseToConn(conn net.Conn, res Response) error {
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
