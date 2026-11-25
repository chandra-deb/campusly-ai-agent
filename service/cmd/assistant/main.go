package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"campusly.local/assistant/internal/personal"
)

func main() {
	key := os.Getenv("CAMPUSLY_AI_KEY")
	if len(key) < 32 {
		log.Fatal("CAMPUSLY_AI_KEY must be a dedicated key of at least 32 bytes")
	}
	port := os.Getenv("CAMPUSLY_AI_PORT")
	if port == "" {
		port = "8765"
	}
	server := &http.Server{
		Addr:              net.JoinHostPort("127.0.0.1", port),
		Handler:           personal.Handler([]byte(key), personal.Fake{}),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      7 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    8 << 10,
	}
	log.Fatal(server.ListenAndServe())
}
