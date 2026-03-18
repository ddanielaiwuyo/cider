package server

import (
	pb "github.com/persona-mp3/protocols/gen"
)

func handleGameMessage(mgr *Manager, packet *pb.Packet) {
	infoLogger.Println("handling game packet")
	mgr.game <- packet
}

func handleChatMessage(mgr *Manager, msg *pb.ChatMessage) {
	infoLogger.Printf("handling chat msg: %+v\n", msg)
}

func handleUnidentifiedPacket(mgr *Manager, msg *pb.Packet) {
	infoLogger.Printf("handling unidentified packet: %+v\n", msg)
}
