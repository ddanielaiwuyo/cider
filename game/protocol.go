package game

import "net"

const maxId = 2000

type Request struct {
	id   int
	move string
}

type User struct {
	id   int
	conn net.Conn
}
