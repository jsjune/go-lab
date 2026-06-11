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

	// 2. 비동기 구독 (Subscribe)
	// 'orders.*.new' 주소로 발행되는 메시지를 수신하기 위한 구독을 생성합니다.
	// Go NATS 클라이언트에서는 콜백 함수를 인자로 넘겨 메시지가 올 때마다 비동기적으로 처리하도록 합니다.
	orderHandler := func(msg *nats.Msg) {
		fmt.Printf("Received request: %s : %s\n", msg.Subject, string(msg.Data))
	}
	sub, err := nc.Subscribe("orders.*.new", orderHandler)
	if err != nil {
		log.Fatalf("구독 설정 실패: %v", err)
	}
	defer sub.Unsubscribe()

	// NATS 서버에 구독 등록이 완료될 때까지 잠시 대기 (선택사항)
	time.Sleep(100 * time.Millisecond)

	// 3. 메시지 발행 (Publish)
	// Go NATS 클라이언트의 Publish 함수는 []byte 타입을 바디 데이터로 받으므로 별도의 인코더가 필요하지 않습니다.
	_ = nc.Publish("orders.test.new", []byte("테스트01 신규 주문"))
	if err := nc.Publish("orders.test.new", []byte("테스트02 신규 주문")); err != nil {
		log.Printf("발행 실패 (test): %v", err)
	}
	if err := nc.Publish("orders.kr.new", []byte("한국 신규 주문")); err != nil {
		log.Printf("발행 실패 (kr): %v", err)
	}
	if err := nc.Publish("orders.eu.new", []byte("유럽 신규 주문")); err != nil {
		log.Printf("발행 실패 (eu): %v", err)
	}
	if err := nc.Publish("orders.us.new", []byte("미국 신규 주문")); err != nil {
		log.Printf("발행 실패 (us): %v", err)
	}

	// 4. 비동기 수신 메시지가 출력될 때까지 잠시 대기 후 안전하게 연결 종료 (Drain)
	time.Sleep(1 * time.Second)

	fmt.Println("NATS 연결을 안전하게 종료(Drain)합니다...")
	if err := nc.Drain(); err != nil {
		log.Printf("Drain 실패: %v", err)
	}
}
