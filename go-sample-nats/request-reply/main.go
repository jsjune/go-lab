package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	message1 := func(msg *nats.Msg) {
		log.Printf("Received a message: %s: %s", msg.Subject, string(msg.Data))

		responseObj := UserResponse{
			Name:  "John",
			Email: "john@example.com",
		}

		responseBytes, err := json.Marshal(responseObj)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			return
		}

		if err := msg.Respond(responseBytes); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
	_, err = nc.Subscribe("user.get", message1)
	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
	}

	message2 := func(msg *nats.Msg) {
		log.Printf("Received a message 22: %s: %s", msg.Subject, string(msg.Data))

		responseObj := UserResponse{
			Name:  "John",
			Email: "john@example.com",
		}

		responseBytes, err := json.Marshal(responseObj)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			return
		}

		if err := msg.Respond(responseBytes); err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}

	_, err = nc.Subscribe("user.get", message2)
	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
	}

	res, err := nc.Request("user.get", []byte("Requesting user data"), nats.DefaultTimeout)
	if err != nil {
		log.Printf("Error making request: %v", err)
	} else {
		log.Printf("Received response: %s", string(res.Data))
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("NATS 연결을 안전하게 종료(Drain)합니다...")
	if err := nc.Drain(); err != nil {
		log.Printf("Drain 실패: %v", err)
	}
}
