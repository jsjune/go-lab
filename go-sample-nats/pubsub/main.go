package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// 1. Nats 서버 연결
	nc, err := nats.Connect("http://localhost:4222")
	if err != nil {
		log.Fatalf("NATS 연결 실패: %v", err)
	}
	defer nc.Close()

	handler := func(msg *nats.Msg) {
		log.Printf("Received message on subject %s: %s", msg.Subject, string(msg.Data))
	}
	_, err = nc.Subscribe("orders.created", handler)
	if err != nil {
		log.Fatalf("구독 설정 실패: %v", err)
	}

	orderData := []byte("새로운 주문이 생성되었습니다!")
	if err := nc.Publish("orders.created", orderData); err != nil {
		log.Printf("발행 실패: %v", err)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("NATS 연결을 안전하게 종료(Drain)합니다...")
	if err := nc.Drain(); err != nil {
		log.Printf("Drain 실패: %v", err)
	}
}
