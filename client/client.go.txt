package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"zod/protocol"
)

func sendRequest(req protocol.Request, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	content, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Write(content); err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Printf(" #: %s\n", buffer[:n])
	}
}

func run() []protocol.Request {
	req := protocol.Request{
		UserId: 1,
		Msg:    "This should doth has work??\n",
	}
	return []protocol.Request{req}
}

func main() {
	var wg sync.WaitGroup
	for _, req := range run() {
		wg.Add(1)
		go sendRequest(req, &wg)
	}

	wg.Wait()
}
