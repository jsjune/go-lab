package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

const (
	defaultNATSURL = nats.DefaultURL
	defaultHTTP    = ":8080"
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

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/stream", streamHandler(nc))

	addr := env("SERVER_A_ADDR", defaultHTTP)
	log.Printf("server-a listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server-a: %v", err)
	}
}

func streamHandler(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.DefaultQuery("clientId", "demo-client")
		inbox := nats.NewInbox()

		msgCh := make(chan *nats.Msg, 8)
		sub, err := nc.ChanSubscribe(inbox, msgCh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer func() {
			if err := sub.Unsubscribe(); err != nil {
				log.Printf("unsubscribe inbox: %v", err)
			}
		}()

		reqBytes, err := json.Marshal(streamRequest{ClientID: clientID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := nc.PublishRequest(streamSubject, inbox, reqBytes); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := nc.Flush(); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		headers := c.Writer.Header()
		headers.Set("Content-Type", "text/event-stream")
		headers.Set("Cache-Control", "no-cache")
		headers.Set("Connection", "keep-alive")
		headers.Set("X-Accel-Buffering", "no")

		timeout := time.NewTimer(5 * time.Second)
		defer timeout.Stop()

		c.Stream(func(w io.Writer) bool {
			select {
			case <-c.Request.Context().Done():
				return false
			case <-timeout.C:
				writeSSE(c, "error", gin.H{"error": "stream timeout"})
				return false
			case msg := <-msgCh:
				var event streamEvent
				if err := json.Unmarshal(msg.Data, &event); err != nil {
					writeSSE(c, "error", gin.H{"error": err.Error()})
					return false
				}

				writeSSE(c, "message", event)
				return !event.Done
			}
		})
	}
}

func writeSSE(c *gin.Context, event string, data any) {
	c.SSEvent(event, data)
	c.Writer.Flush()
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
