package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getServer(bindAddress string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/", handler)
	server := &http.Server{
		Addr:    bindAddress,
		Handler: router,
	}
	return server
}
