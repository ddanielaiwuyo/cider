package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
)

func readFromStdin(ctx context.Context, in chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		case in <- scanner.Text():
		}
	}
}

func readFromServer(ctx context.Context, conn net.Conn, out chan<- []byte) error {
	buff := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := conn.Read(buff)
			if err != nil {
				return fmt.Errorf(" read-err: %w", err)
			}

			out <- buff[:n]

		}
	}
}
