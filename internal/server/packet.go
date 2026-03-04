package server

import (
	"encoding/binary"
	"fmt"

	pb "github.com/persona-mp3/protocols/github.com/persona-mp3/protocols"
	"google.golang.org/protobuf/proto"
)

// Converts extractedPacketData to a pb.Packet defined
// according to spec
func parsePacketData(data []byte) (*pb.Packet, error) {
	msg := &pb.Packet{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("could not parse data packet %w", err)
	}
	return msg, nil
}

func createPaintPacket(dest int32) ([]byte, error) {
	msg := createPaintMessage()
	packet := pb.Packet{
		From: serverId,
		Dest: dest,
		Payload: &pb.Packet_Paint{
			Paint: msg,
		},
	}

	wirePacket, err := MarshallPacket(&packet)
	if err != nil {
		return []byte{}, err
	}
	return wirePacket, nil
}


func MarshallPacket(packet *pb.Packet) ([]byte, error) {
	data, err := proto.Marshal(packet)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshall packet: %w , %+v", err, packet)
	}

	header := make([]byte, headerLength)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	wirePacket := append(header, data...)
	return wirePacket, nil
}
