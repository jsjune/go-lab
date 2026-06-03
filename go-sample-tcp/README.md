# go-sample-tcp

Go 표준 라이브러리(`net`)만으로 구현한 TCP 채팅 서버 샘플입니다.
여러 클라이언트가 동시에 접속해 메시지를 주고받습니다.

---

## 폴더 구조

```
go-sample-tcp/
├── cmd/
│   ├── server/
│   │   └── main.go     # TCP 서버: hub 관리 + 연결 수락 루프
│   └── client/
│       └── main.go     # TCP 클라이언트: 입출력 고루틴 분리
└── go.mod
```

---

## 실행

터미널을 여러 개 열어 서버를 먼저 시작하고, 클라이언트를 여러 개 접속시킵니다.

```bash
# 터미널 1 — 서버 시작
go run ./cmd/server

# 터미널 2 — 클라이언트 접속
go run ./cmd/client

# 터미널 3 — 클라이언트 추가 접속
go run ./cmd/client
```

접속 후 닉네임을 입력하면 채팅이 시작됩니다. `/quit`을 입력하면 종료합니다.

---

## 학습 포인트

| 코드 | 내용 |
|------|------|
| `hub.clients map[net.Conn]string` | 연결된 클라이언트를 맵으로 관리 (conn → 닉네임) |
| `hub.mu sync.Mutex` | 여러 고루틴이 동시에 clients 맵 수정 시 경쟁 조건 방지 |
| `go h.handle(conn)` | 클라이언트마다 고루틴 하나 — OS 스레드가 아닌 경량 고루틴 |
| `bufio.Scanner` | TCP 스트림을 줄 단위로 읽기 |
| `hub.broadcast` | 발신자를 제외한 모든 연결에 메시지 전송 |
| `client: go io.Copy(os.Stdout, conn)` | 서버 수신을 별도 고루틴으로 분리 — 송수신 동시 처리 |
| `defer h.leave(conn)` | 고루틴 종료 시(접속 끊김) 자동으로 클라이언트 제거 |

---

## 아키텍처

```
[Client A] ─────────────────────────────────────────────────┐
[Client B] ──── go h.handle(conn) ──── hub ──── broadcast ──┤
[Client C] ─────────────────────────────────────────────────┘
                     ↑
              클라이언트마다 고루틴 하나
              hub는 sync.Mutex로 보호
```

### 서버 흐름

```
net.Listen(":9000")
  └── for { conn := ln.Accept(); go h.handle(conn) }
            │
            ├── 닉네임 입력 받기
            ├── h.join(conn, name)  → 입장 브로드캐스트
            ├── for scanner.Scan()  → 메시지 브로드캐스트
            └── defer h.leave(conn) → 퇴장 브로드캐스트
```

### 클라이언트 흐름

```
net.Dial("localhost:9000")
  ├── go io.Copy(os.Stdout, conn)  → 서버 메시지 수신 (별도 고루틴)
  └── bufio.Scanner(os.Stdin)      → 키보드 입력 → 서버 전송 (메인 고루틴)
```

---

## 명령어

| 입력 | 동작 |
|------|------|
| 텍스트 입력 후 Enter | 접속 중인 다른 클라이언트 전체에 브로드캐스트 |
| `/quit` | 접속 종료 |
