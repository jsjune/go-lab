package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// hub는 연결된 모든 클라이언트를 관리하고 브로드캐스트를 담당
type hub struct {
	mu      sync.Mutex
	clients map[net.Conn]string // conn → 닉네임
}

func newHub() *hub {
	return &hub{clients: make(map[net.Conn]string)}
}

func (h *hub) join(conn net.Conn, name string) {
	h.mu.Lock()
	h.clients[conn] = name
	h.mu.Unlock()
	h.broadcast(conn, fmt.Sprintf(">>> %s 님이 입장했습니다.", name))
}

func (h *hub) leave(conn net.Conn) {
	h.mu.Lock()
	name := h.clients[conn]
	delete(h.clients, conn)
	h.mu.Unlock()
	h.broadcast(conn, fmt.Sprintf("<<< %s 님이 퇴장했습니다.", name))
}

// 보낸 사람을 제외한 모든 클라이언트에게 전송
func (h *hub) broadcast(sender net.Conn, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	line := fmt.Sprintf("[%s] %s\n", time.Now().Format("15:04:05"), msg)
	for conn := range h.clients {
		if conn != sender {
			fmt.Fprint(conn, line)
		}
	}
	log.Print(strings.TrimSpace(line))
}

func (h *hub) handle(conn net.Conn) {
	defer conn.Close()

	fmt.Fprint(conn, "닉네임을 입력하세요: ")
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = conn.RemoteAddr().String()
	}

	h.join(conn, name)
	defer h.leave(conn)

	fmt.Fprintf(conn, "접속 완료. 메시지를 입력하세요 (종료: /quit)\n")

	for scanner.Scan() {
		msg := strings.TrimSpace(scanner.Text())
		if msg == "/quit" {
			break
		}
		if msg == "" {
			continue
		}
		h.broadcast(conn, fmt.Sprintf("%s: %s", name, msg))
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP 채팅 서버 시작 — :9000")

	h := newHub()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		// 클라이언트마다 고루틴 하나 — 스레드 대신 경량 고루틴 사용
		go h.handle(conn)
	}
}
