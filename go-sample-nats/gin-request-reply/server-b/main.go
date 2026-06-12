package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	defaultNATSURL = nats.DefaultURL
	streamSubject  = "sample.stream.request"
)

type streamRequest struct {
	ClientID string `json:"clientId"`
}

type streamEvent struct {
	Index     int    `json:"index"`
	Total     int    `json:"total"`
	Message   string `json:"message"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}

func main() {
	nc, err := nats.Connect(env("NATS_URL", defaultNATSURL))
	if err != nil {
		log.Fatalf("connect nats: %v", err)
	}
	defer nc.Close()

	_, err = nc.Subscribe(streamSubject, func(msg *nats.Msg) {
		if msg.Reply == "" {
			log.Printf("skip request without reply subject")
			return
		}

		var req streamRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			log.Printf("decode request: %v", err)
			return
		}

		log.Printf("stream request received: clientId=%s reply=%s", req.ClientID, msg.Reply)
		for i := 1; i <= 5; i++ {
			event := streamEvent{
				Index:     i,
				Total:     5,
				Message:   fmt.Sprintf("mock stream data %d for %s", i, req.ClientID),
				Done:      i == 5,
				CreatedAt: time.Now().Format(time.RFC3339Nano),
			}

			eventBytes, err := json.Marshal(event)
			if err != nil {
				log.Printf("encode event: %v", err)
				return
			}
			if err := nc.Publish(msg.Reply, eventBytes); err != nil {
				log.Printf("publish reply: %v", err)
				return
			}
			if err := nc.Flush(); err != nil {
				log.Printf("flush reply: %v", err)
				return
			}

			time.Sleep(300 * time.Millisecond)
		}
	})
	if err != nil {
		log.Fatalf("subscribe %s: %v", streamSubject, err)
	}

	if err := nc.Flush(); err != nil {
		log.Fatalf("flush subscription: %v", err)
	}

	log.Printf("server-b subscribed to %s", streamSubject)
	select {}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
