# go-sample-nats

NATS 공식 Go 클라이언트(`github.com/nats-io/nats.go`)로 기본 Pub/Sub와 와일드카드 라우팅을 연습하는 샘플입니다.

---

## 폴더 구조

```text
go-sample-nats/
|-- docker/
|   |-- docker-compose.yml   # NATS 서버와 NUI 웹 대시보드 실행
|   `-- nats.conf            # NATS 서버 설정
|-- pubsub/
|   `-- main.go              # 기본 Pub/Sub 예제
|-- request-reply/
|   `-- main.go              # Request/Reply 예제
|-- routing/
|   `-- main.go              # 와일드카드 subject 라우팅 예제
|-- go.mod
`-- go.sum
```

---

## 사전 준비

- Go 1.26.1 이상
- Docker Desktop 또는 Docker Compose 실행 환경

의존성은 `go.mod`에 선언되어 있습니다.

```bash
github.com/nats-io/nats.go
```

---

## NATS 서버 실행

`go-sample-nats/docker` 디렉터리에서 NATS와 NUI를 실행합니다.

```bash
cd go-sample-nats/docker
docker compose up -d
```

접속 정보:

| 용도 | URL / 포트 |
|---|---|
| NATS 클라이언트 연결 | `localhost:4222` |
| NATS HTTP 모니터링 | `http://localhost:8222` |
| NUI 웹 대시보드 | `http://localhost:31311` |

NUI에서 모니터링을 연결할 때는 Metrics Source를 `HTTP`로 두고 URL에 `http://nats-server:8222`를 입력합니다. NUI 컨테이너와 NATS 컨테이너가 같은 Docker 네트워크에 있으므로 컨테이너 이름인 `nats-server`를 사용합니다.

서버 중지:

```bash
cd go-sample-nats/docker
docker compose down
```

---

## Go 예제 실행

`go-sample-nats` 디렉터리에서 실행합니다.

```bash
cd go-sample-nats
```

기본 Pub/Sub 예제:

```bash
go run ./pubsub
```

Request/Reply 예제:

```bash
go run ./request-reply
```

와일드카드 라우팅 예제:

```bash
go run ./routing
```

---

## 예제 설명

### `pubsub/main.go`

가장 기본적인 Pub/Sub 예제입니다.

1. `nats.Connect("http://localhost:4222")`로 NATS 서버에 연결합니다.
2. `nc.Subscribe("orders.created", handler)`로 `orders.created` subject를 구독합니다.
3. `nc.Publish("orders.created", []byte(...))`로 메시지를 발행합니다.
4. 콜백 핸들러가 발행된 메시지를 비동기로 수신해 로그로 출력합니다.
5. `nc.Drain()`으로 남은 메시지를 처리한 뒤 연결을 정리합니다.

핵심 코드:

```go
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
```

### `request-reply/main.go`

NATS의 Request/Reply 패턴 예제입니다.

1. `nc.Subscribe("user.get", handler)`로 요청을 받을 구독자를 등록합니다.
2. 요청 메시지를 받으면 `msg.Respond(...)`로 응답을 보냅니다.
3. `nc.Request("user.get", []byte(...), nats.DefaultTimeout)`으로 요청을 보내고 응답을 기다립니다.

핵심 코드:

```go
_, err = nc.Subscribe("user.get", func(msg *nats.Msg) {
    responseBytes, err := json.Marshal(UserResponse{
        Name:  "John",
        Email: "john@example.com",
    })
    if err != nil {
        log.Printf("Error marshaling response: %v", err)
        return
    }

    if err := msg.Respond(responseBytes); err != nil {
        log.Printf("Error sending response: %v", err)
    }
})

res, err := nc.Request("user.get", []byte("Requesting user data"), nats.DefaultTimeout)
```

### `routing/main.go`

NATS subject 와일드카드 라우팅 예제입니다.

```go
sub, err := nc.Subscribe("orders.*.new", orderHandler)
```

`orders.*.new`는 가운데 토큰 하나가 어떤 값이든 매칭합니다.

매칭되는 subject:

```text
orders.test.new
orders.kr.new
orders.eu.new
orders.us.new
```

매칭되지 않는 subject:

```text
orders.new
orders.us.east.new
```

---

## 학습 포인트

| 코드 / 개념 | 설명 |
|---|---|
| `nats.Connect(...)` | NATS 서버와 클라이언트 연결을 생성합니다. |
| `defer nc.Close()` | `main()`이 끝날 때 연결을 닫도록 예약합니다. Java의 `finally` 또는 try-with-resources 정리 동작과 비슷합니다. |
| `nc.Subscribe(subject, handler)` | subject를 구독하고 메시지가 오면 콜백 함수를 비동기로 실행합니다. |
| `*nats.Msg` | 수신한 메시지 객체입니다. `Subject`, `Data` 등을 가집니다. |
| `nc.Publish(subject, []byte(...))` | 특정 subject로 메시지를 발행합니다. 메시지 본문은 `[]byte`입니다. |
| `time.Sleep(...)` | 예제 프로그램이 바로 종료되지 않도록 잠시 대기합니다. 실제 서버 코드에서는 보통 다른 방식으로 생명주기를 관리합니다. |
| `nc.Drain()` | 남은 메시지 처리를 기다린 뒤 안전하게 연결을 종료합니다. |
| `sub.Unsubscribe()` | 구독을 해제합니다. |

---

## NATS Subject 와일드카드

NATS subject는 점(`.`)으로 토큰을 나눕니다.

### 단일 토큰 와일드카드: `*`

정확히 하나의 토큰을 대체합니다.

```text
orders.*.new
```

매칭:

```text
orders.kr.new
orders.us.new
```

불일치:

```text
orders.new
orders.us.east.new
```

### 다중 토큰 와일드카드: `>`

subject의 마지막에만 사용할 수 있고, 이후의 하나 이상의 토큰을 대체합니다.

```text
orders.>
```

매칭:

```text
orders.new
orders.us.new
orders.us.east.new
```

불일치:

```text
order.new
```

---

## 참고

이 예제에서는 학습 편의를 위해 `log.Fatalf`를 사용합니다. 단, `log.Fatal`/`log.Fatalf`는 내부적으로 `os.Exit(1)`을 호출하므로 `defer`가 실행되지 않습니다.

실무 코드에서는 보통 다음처럼 `run()` 함수에서 에러를 반환하고, `main()`에서 한 번만 종료 처리합니다.

```go
func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    nc, err := nats.Connect("http://localhost:4222")
    if err != nil {
        return err
    }
    defer nc.Close()

    return nil
}
```
