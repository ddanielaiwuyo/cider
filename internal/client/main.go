package main

import (
	"log"

	"github.com/persona-mp3/internal/client/impl"
)

func main() {

	log.Println("runnning client package")
	impl.DialServer(4000,
		impl.AuthCredentials{Username: "rick"},
	) // use args instead
}
