package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal("서버 연결 실패:", err)
	}
	defer conn.Close()

	// 서버 → 터미널 출력 (별도 고루틴으로 수신과 송신을 동시에 처리)
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Println("\n서버 연결 종료")
			os.Exit(0)
		}
	}()

	// 터미널 입력 → 서버 전송
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintln(conn, scanner.Text())
	}
}
