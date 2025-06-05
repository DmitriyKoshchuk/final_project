package server

import (
	"log"
	"net/http"

	"go1f/pkg/api"
)

func Run() error {
	mux := http.NewServeMux()
	api.Init(mux)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	log.Println("Starting server at :7540")
	return http.ListenAndServe(":7540", mux)
}
