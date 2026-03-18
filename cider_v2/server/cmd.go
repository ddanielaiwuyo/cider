package server

import (
	pb "github.com/persona-mp3/protocols/gen"
)

type InternalCommand int

const (
	Deliver InternalCommand = iota
)

var gameServerId = 01

type Command struct {
	Id      int
	CmdType InternalCommand
	Packet  *pb.Packet
}

func deliverCommand(packet *pb.Packet) *Command {
	return &Command{
		Id:      gameServerId,
		CmdType: Deliver,
		Packet:  packet,
	}
}
