package main

import (
	"log"
	"zod/game/chess"
)

func main() {
	if err := chess.Run(); err != nil {
		log.Fatal(err)
	}
}
