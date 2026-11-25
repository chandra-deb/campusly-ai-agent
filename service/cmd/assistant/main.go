package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"campusly.local/assistant/internal/personal"
)

func main() {
	smokeKeyFile := flag.String("cerebras-smoke-key-file", "", "enable one Cerebras request using this owner-only key file")
	flag.Parse()
	key := os.Getenv("CAMPUSLY_AI_KEY")
	if len(key) < 32 {
		log.Fatal("CAMPUSLY_AI_KEY must be a dedicated key of at least 32 bytes")
	}
	port := os.Getenv("CAMPUSLY_AI_PORT")
	if port == "" {
		port = "8765"
	}
	handler := personal.Handler([]byte(key), personal.Fake{})
	if *smokeKeyFile != "" {
		info, err := os.Lstat(*smokeKeyFile)
		if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 || info.Size() > 4096 {
			log.Fatal("Cerebras smoke needs an owner-only regular key file")
		}
		data, err := os.ReadFile(*smokeKeyFile)
		providerKey := strings.TrimSpace(string(data))
		if err != nil || providerKey == "" || strings.ContainsAny(providerKey, "\r\n") {
			log.Fatal("Cerebras smoke needs a key")
		}
		hosted := personal.Handler([]byte(key), personal.NewCerebrasSmoke(providerKey))
		fake := handler
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == personal.ContextPath {
				hosted.ServeHTTP(w, r)
				return
			}
			fake.ServeHTTP(w, r)
		})
	}
	server := &http.Server{
		Addr:              net.JoinHostPort("127.0.0.1", port),
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      7 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    8 << 10,
	}
	log.Fatal(server.ListenAndServe())
}
