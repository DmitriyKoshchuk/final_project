package main

import (
	"log"

	"github.com/DmitriyKoshchuk/final_project/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
