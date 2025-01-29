package main

import (
	"log"

	"github.com/coronellw/go-microservices/internal/database"
	"github.com/coronellw/go-microservices/internal/server"
)

func main() {
	db, err := database.NewDatabaseClient()

	if err != nil {
		log.Fatalf("failed to initialize Database Client: %s", err)
	}

	srv := server.NewEchoServer(db.(database.Client))
	if err := srv.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
