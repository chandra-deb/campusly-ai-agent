package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"campusly.local/assistant/internal/personal"
)

func main() {
	port := os.Getenv("CAMPUSLY_AI_PORT")
	if port == "" {
		port = "8765"
	}
	key := os.Getenv("CAMPUSLY_AI_KEY")
	payload, err := json.Marshal(personal.Request{
		Version: 1, Scope: personal.Scope{Actor: 11, Organization: 21},
		Expires: time.Now().Add(5 * time.Second).UnixMilli(),
		Prompt:  "Give me a workshop checklist",
	})
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://"+net.JoinHostPort("127.0.0.1", port)+personal.Path, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set(personal.Signature, personal.Sign([]byte(key), payload))
	client := http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{}}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("fixture returned HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<10))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(body))
}
