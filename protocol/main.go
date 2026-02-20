package protocol

import (
	"fmt"
	"strings"
)

type MessageType int

const (
	ConnectTo MessageType = iota
	Unknown
)

func ParseMessage(msg []byte) (MessageType, string) {
	str_fmt := fmt.Sprintf("%s", msg)

	switch true {
	case strings.Contains(str_fmt, "connect-to"):
		_, addr, _ := strings.Cut(str_fmt, "connect-to")
		addr = strings.TrimSpace(addr)
		return ConnectTo, addr
	}

	return Unknown, str_fmt
}
