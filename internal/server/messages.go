package server

import (
	"fmt"
	"strings"
)

func createPaintMessage() string {
	var msg strings.Builder
	msg.WriteString("Connected Users;")

	for id := range connectedUsers {
		fmt.Fprintf(&msg, "User %d;", id)
	}

	return msg.String()
}
