package server

import (
	"log"

	"github.com/google/uuid"
	pb "github.com/persona-mp3/protocols/github.com/persona-mp3/protocols"
)

func createPaintMessage() *pb.PaintMessage {
	// should just get a list of all
	// connected users, possibly provided from?
	activeUsers := []*pb.User{
		{Username: "gopls"},
		{Username: "are_you_ladies_man_217?"},
		{Username: "a_blow_fish!"},
	}

	paintMsg := pb.PaintMessage{
		ConnectedUsers: activeUsers,
		OneTimeId:      uuid.New().String(),
	}
	log.Println(" [debug] created id for user")
	return &paintMsg
}
