package impl

import (
	"log"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	// "os/signal"

	pb "github.com/persona-mp3/protocols/gen"

	pack "github.com/persona-mp3/internal/packet"
)

func fromServer(ctx context.Context, conn net.Conn) <-chan *pb.Packet {
	response := make(chan *pb.Packet, 16)
	go func() {
		// defer close(response)
		for {
			content, err := pack.ReadWirePacket(conn, headerSize)
			if err != nil {
				if errors.Is(err, io.EOF) {
					slog.Error("server has disconnected!", "err", err)
					return
				} else {
					slog.Error("unexpected error", "err", err)
					return
				}
			}

			packet, err := pack.ParseWirePacket(content)
			if err != nil {
				slog.Error("error in parsing wire packet", "err", err)
				continue
			}


			select {
			case <-ctx.Done():
				slog.Info(" ctx called!")
				return
			case response <- packet:
			}
		}
	}()
	return response
}

func fromStdin(ctx context.Context) <-chan string {
	// ct, cancel := context.WithCancel(ctx)
	// _ = ct
	// _c=Vkk
	stdin := make(chan string)
	in := os.Stdin
	scanner := bufio.NewScanner(in)
	// c := make(chan os.Signal, 1)
	// signal.Notify(c)
	go func() {
		defer func() {
			fmt.Println("someone returned here!")
			fmt.Printf("stdin status: %+v, fd: %+v\n", in, in.Fd())
			if in == nil {
				fmt.Println("stdin is nil!!!!", in)
			} else {
				fmt.Println("stdin is still open???!!!!", in)
			}

			// fmt.Println("err from scanner?", scanner.Err().Error())
			if scanner == nil {
				fmt.Println("we caught the culprit", scanner)
			} else {
				fmt.Println("scaner open??", scanner.Err())
			}

			// in.Close()

			fmt.Printf("closed -> %+v\n", in)
			fmt.Printf("ctx err: %s\n", ctx.Err())
			if len(scanner.Text()) != 0 {
				print("non e o yo", scanner.Text())
				// close(stdin)
			} else {
				fmt.Printf("after defer stll active %s, %d\n", scanner.Text(), len(scanner.Text()))
				// close(stdin)
			}
			close(stdin)

		}()
		// defer close(stdin)
		for scanner.Scan() {
			// fmt.Printt(" [*] ")
			if err := scanner.Err(); err != nil {
				fmt.Println("scanner error occured", err)
				panic("paniced")
			}
			if in == nil {
				panic("stdin const const const cosnt")
			}
			// fmt.Println("moding")
			value := scanner.Text()
			log.Println("[debug] scanner ret valu ->", value)
			select {
			case <-ctx.Done():
				fmt.Println(" [debug] ctx called in scanner", ctx.Err().Error())
				return

			case stdin <- value:
				fmt.Print(" [*] ")
				// default:
				// case sig := <-c:
				// 	fmt.Printf(" [debug] signal recvd: %+v\n", sig)
				// 	if sig == os.Interrupt {
				// 		panic("deadly sig")
				// 	}
			}
		}
	}()

	slog.Info("connected to stdin successfully")
	return stdin
}

func toServer(ctx context.Context, conn net.Conn) chan<- *pb.Packet {
	writer := make(chan *pb.Packet)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case packet, ok := <-writer:
				fmt.Println(" [debug]watioting fro writer   dvknkvkjndskjvnskj")
				if !ok {
					slog.Error("writer channel closed", "isopen", ok)
					return
				}
				content, err := pack.MarshallPacket(packet, headerSize)
				if err != nil {
					slog.Error("while marshalling", "err", err)
					continue
				}

				if _, err := conn.Write(content); err != nil {
					slog.Error("could not write message to sever", "err", err)
					return
				}
			}
		}
	}()
	return writer
}
