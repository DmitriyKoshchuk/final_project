package main

import (
	"log"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
