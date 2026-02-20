package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn
	addr string
	Hub  map[net.Conn]bool
}

func CreatePeer(conn net.Conn) Peer {
	return Peer{
		conn: conn,
		addr: conn.RemoteAddr().String(),
		Hub:  make(map[net.Conn]bool),
	}
}

func validateTCPAddr(addr string) error {
	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not parse addr %w", err)
	}
	return nil
}

func (peer Peer) ConnectTo(destAddr string) error {
	// check the connection pool if they exists
	if err := validateTCPAddr(destAddr); err != nil {
		return err
	}
	foundDest := false
	var destConn net.Conn
	for conn := range connectionPool {
		if conn.RemoteAddr().String() == destAddr &&
			conn.RemoteAddr().String() != peer.addr {
			foundDest = true
			destConn = conn
			peer.Hub[destConn] = true
		}
	}

	if !foundDest {
		return fmt.Errorf(" dest could not be found")
	}

	sendConfirmationMsg(peer.conn, destConn)
	buff := CreateBuffer(1024)
	go func() {
		for {
			// read from originator and read from dest
			n, err := peer.conn.Read(buff)
			if err != nil && errors.Is(err, io.EOF) {
				slog.Info(" client disconnected:", "addr:", peer.addr)
				return
			} else if err != nil && !errors.Is(err, io.EOF) {
				slog.Error(" read_err:", "error", err)
				peer.conn.Close()
				return
			}

			if _, err := destConn.Write(buff[:n]); err != nil {
				slog.Error(" could not write to dest, ", "err", err)
				destConn.Close()
				return
			}

			// just breathe, we'll refactor this later
			// to use channels
			go func() {
				n, err := destConn.Read(buff)
				if err != nil && errors.Is(err, io.EOF) {
					slog.Info(" client disconnected:", "addr:", peer.addr)
					return
				} else if err != nil && !errors.Is(err, io.EOF) {
					slog.Error(" read_err:", "error", err)
					peer.conn.Close()
					return
				}

				if _, err := peer.conn.Write(buff[:n]); err != nil {
					slog.Error(" could not write to dest, ", "err", err)
					destConn.Close()
					return
				}
			}()
		}
	}()

	return nil
}

func CreateBuffer(size uint16) []byte {
	return make([]byte, size)
}

func closer(conn net.Conn) {
	conn.Close()
}

func sendConfirmationMsg(from, to net.Conn) error {
	fmt.Fprintf(to, "TO BE INTERLINKED?: %s\n", from.LocalAddr().String())
	fmt.Fprintf(from, "INTERLINKED with %s\n", to.LocalAddr().String())
	return nil
}

// _, err := from.Write([]byte(interlinkedMsgToOrigin))
// if err != nil {
// 	log.Println("couldn't update originator", err)
// 	to.Write([]byte("yeah bro... wanted to interlink? senpoku\n"))
// 	closer(from)
// 	return err
// }

// if _, err := to.Write([]byte(fmt.Sprintf("TO BE INTERLINKED?: %s\n", from.LocalAddr().String()))); err != nil {
// 	log.Println("couldn't update destinator", err)
// 	closer(to)
// 	from.Write([]byte("yeah bro... could not interlink? senpoku"))
// 	return err
// }
// interlinkedMsgToOrigin := fmt.Sprintf("INTERLINKED with %s\n", to.LocalAddr().String())
