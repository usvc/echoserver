package main

import (
	"net/http"
)

func getServer(bindAddress string) *http.Server {
	server := &http.Server{
		Addr:    bindAddress,
		Handler: http.HandlerFunc(handler),
	}
	return server
}
