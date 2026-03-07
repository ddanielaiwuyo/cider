package server

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	pb "github.com/persona-mp3/protocols/gen"
)

func CreateGameNewGameSession(mgr *manager, req *pb.NewGameMessage) {
	slog.Info("create new game session", "for", req.From, "with", req.Dest)
	gameSessionId := uuid.NewString()
	homePlayer, connected := mgr.connections[connId(req.From)]
	if !connected {
		slog.Info("cannot continue creating game sessions, home player not found")
		return
	}

	awayPlayer, isOnline := mgr.connections[connId(req.Dest)]
	if !isOnline {
		slog.Info("cannot create game session, awayPlayer isn't connected")
		return
	}

	player1, player2 := newPlayer(&homePlayer), newPlayer(&awayPlayer)

	currSession := &GameSession{
		SessionId: gameSessionId,
		Players:   []*Player{player1, player2},
	}

	mgr.GameSessions[gameSessionId] = currSession

	info := fmt.Sprintf(`
	Starting Game
	Challenger: %s
	Away: %s
	`, homePlayer.username, awayPlayer.username)
	mgr.deliver <- &pb.Packet{
		From: "server",
		Dest: req.From,
		Payload: &pb.Packet_NewGameRes{
			NewGameRes: &pb.NewGameResponse{
				Ssid: gameSessionId,
				Info: &info,
				From: "server",
			},
		},
	}
	mgr.deliver <- &pb.Packet{
		From: "server",
		Dest: req.Dest,
		Payload: &pb.Packet_NewGameRes{
			NewGameRes: &pb.NewGameResponse{
				Ssid: gameSessionId,
				Info: &info,
				From: "server",
			},
		},
	}

}

func newPlayer(c *Client) *Player {
	return &Player{
		client: c,
		Play:   make(chan string),
	}
}

// Not sure yet, what to tell the client when they provide a non-existent ssid
func HandleGamePacket(mgr *manager, gamePacket *pb.GameMessage) {
	ssid := gamePacket.Ssid
	session, validSession := mgr.GameSessions[ssid]
	if !validSession {
		slog.Info("client sent a non existing game ssid")
		return
	}

	slog.Info("found game session", "id", ssid)
	for _, player := range session.Players {
		fmt.Printf("%+v\n", player)
	}

}
