package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer db.DB.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		db.DB.Close()
		os.Exit(0)
	}()

	err = server.Run()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
