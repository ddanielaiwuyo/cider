package server

import (
	"log"

	pb "github.com/persona-mp3/protocols/gen"
)

func createPaintMessage(mgr *manager, id connId) *pb.PaintMessage {
	// should just get a list of all
	// connected users, possibly provided from the manager
	// activeUsers := []*pb.User{
	// 	{Username: "gopls"},
	// 	{Username: "are_you_ladies_man_217?"},
	// 	{Username: "a_blow_fish!"},
	// }

	var connections = make([]*pb.User, len(mgr.connections))
	for connId, c := range mgr.connections {
		u := &pb.User{
			Username: c.username,
			Id:       string(connId),
		}

		connections = append(connections, u)
	}

	paintMsg := pb.PaintMessage{
		ConnectedUsers: connections,
		OneTimeId:      string(id),
	}
	log.Println(" [debug] created id for user")
	return &paintMsg
}

func createAuthSuccessMessage(code int32, content string) *pb.AuthSuccess {
	defaultContent := "Auth Successfull"
	if len(content) <= 0 {
		content = defaultContent
	}

	return &pb.AuthSuccess{Code: code, Content: content}
}
