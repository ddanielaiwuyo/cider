package main

import (
	"context"
	"fmt"
	"log"
	"os"

	db "github.com/persona-mp3/internal/database"
	"github.com/persona-mp3/internal/server"
)

func main() {
	//  use ENV instead!
	dbConf := &db.DBConfig{
		Username: "persona",
		Password: "persona-mp3",
		Database: "cidervine",
		Port:     5432,
	}

	conn, err := db.Connect(dbConf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	manager := server.NewManager(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go manager.Listen(ctx)

	if err := server.RunServer(manager); err != nil {
		log.Fatal(err)
	}
}
