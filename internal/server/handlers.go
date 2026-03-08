package server

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/google/uuid"
	pb "github.com/persona-mp3/protocols/gen"
)

var defaultTickerRate int32 = 2

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
		Rate:      defaultTickerRate,
	}

	mgr.GameSessions[gameSessionId] = currSession

	info := fmt.Sprintf(`
	  STARTING GAME
	  Challenger: %s
	  Away: %s
	`, homePlayer.username, awayPlayer.username)

	defaultTickerRate := int32(2)
	// for the challenger
	mgr.deliver <- &pb.Packet{
		From: "server",
		Dest: req.From,
		Payload: &pb.Packet_NewGameRes{
			NewGameRes: &pb.NewGameResponse{
				Ssid:       gameSessionId,
				Info:       &info,
				From:       "server",
				Rival:      req.Dest,
				TickerRate: &defaultTickerRate,
			},
		},
	}
	// for the rival
	mgr.deliver <- &pb.Packet{
		From: "server",
		Dest: req.Dest,
		Payload: &pb.Packet_NewGameRes{
			NewGameRes: &pb.NewGameResponse{
				Ssid:       gameSessionId,
				Info:       &info,
				From:       "server",
				Rival:      req.Dest,
				TickerRate: &defaultTickerRate,
			},
		},
	}

}

func newPlayer(c *Client) *Player {
	return &Player{
		client: c,
		// Play:   make(chan string),
	}
}

func HandleGamePacket(mgr *manager, packet *pb.Packet) {
	gameMessage := packet.GetGame()
	ssid := gameMessage.Ssid

	session, validSession := mgr.GameSessions[ssid]
	_ = session

	if !validSession {
		slog.Info("client sent a non existing game ssid")
		return
	}

	for _, player := range session.Players {
		if packet.From == string(player.client.userId) {
			continue
		}
		mgr.deliver <- &pb.Packet{
			From: "server",
			Dest: string(player.client.userId),
			Payload: &pb.Packet_Game{
				Game: &pb.GameMessage{
					Play: gameMessage.Play,
					Ssid: ssid,
				},
			},
		}

		log.Println("[debug] game_play sent to ", player.client.username)
	}

	log.Printf("player %s: %s\n", packet.From, gameMessage.Play)
	// not sure if the clients need to send their rival in
	// every request from now on, because the server automatically handles that
	// by just broadcasting to all the gamers in a session
	// log.Println("[debug] sent to rival:", gameMessage.Rival, packet.From)
	log.Println("[debug] broadcast game to all players")

}

/*
A sends CHALLENGE B
B -> Rival: A
A -> Rival: B
*/
