package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/persona-mp3/shared"
)

func readFromStdin(ctx context.Context) <-chan string {
	log.Printf("[ch] reading from stdin")
	stdin := make(chan string)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		defer close(stdin)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case stdin <- scanner.Text():
				fmt.Print(" [*] ")
			}
		}
	}()

	return stdin
}

func writeToServer(msg string, conn net.Conn) error {
	// we're actually supposed to store the uuid
	// the server will provide on first connection
	// for future auth and to avoid confusion.
	// But for now, we could
	// just hardcode it
	req := shared.Message{
		Dest:        2,
		From:        1,
		MessageType: shared.ChatMessage,
		Content:     msg,
	}

	content, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not parse msg to send %w", err)
	}

	if _, err := io.Copy(conn, bytes.NewReader(content)); err != nil {
		return fmt.Errorf("could not write to server: %w", err)
	}

	return nil
}

func readFromServer(ctx context.Context, conn net.Conn) <-chan shared.Message {
	log.Printf("[ch] reading from server")
	res := make(chan shared.Message)
	go func() {
		buff := make([]byte, 1024)
		defer close(res)
		for {
			select {
			case <-ctx.Done():
				slog.Error("[ch] context done", "", ctx.Err().Error())
				return
			default:

				n, err := conn.Read(buff)
				if err != nil {
					slog.Error("read error from server", "", err)
					return
				}

				dest := make([]byte, n)
				copy(dest, buff[:n])

				var msg shared.Message
				if err := json.Unmarshal(dest, &msg); err != nil {
					log.Println(string(dest))
					slog.Error("could not parse server res", "", err)
					return
				}

				res <- msg
			}
		}
	}()
	return res
}
