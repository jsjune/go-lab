# go-lab

Go 언어 학습용 샘플 프로젝트 모음입니다.
각 프로젝트는 독립된 Go 모듈로, 주제별로 핵심 개념을 집중적으로 다룹니다.

---

## 프로젝트 목록

| 프로젝트 | 기술 스택 | 주제 |
|----------|-----------|------|
| [go-sample-http](./go-sample-http) | Gin, gRPC, SQLite | 레이어드 아키텍처, REST API, gRPC |
| [go-sample-stdlib-http](./go-sample-stdlib-http) | 표준 라이브러리 | `net/http` 직접 구현, 인터페이스 패턴 |
| [go-sample-cli](./go-sample-cli) | Cobra | CLI 도구 개발, JSON 파일 영속화 |
| [go-sample-tcp](./go-sample-tcp) | 표준 라이브러리 | TCP 소켓, 고루틴, 동시성 |

---

## 추천 학습 순서

```
1. go-sample-stdlib-http   표준 라이브러리로 HTTP 기초 이해
        ↓
2. go-sample-http          Gin + gRPC + SQLite로 실전 서버 구조 학습
        ↓
3. go-sample-cli           Cobra로 CLI 도구 개발 패턴 학습
        ↓
4. go-sample-tcp           고루틴과 동시성 직접 구현
```

---

## 각 프로젝트 요약

### [go-sample-http](./go-sample-http)

Gin 프레임워크와 SQLite를 사용한 REST API 서버입니다.
HTTP와 gRPC 두 트랜스포트가 동일한 UseCase를 공유하는 구조를 보여줍니다.

```
cmd/server/  ← HTTP 서버 진입점
cmd/grpc/    ← gRPC 서버 진입점
internal/item/
  entity.go       도메인 구조체
  dto.go          요청/응답 DTO
  repository.go   DB 접근 레이어 (인터페이스 + SQLite 구현)
  usecase.go      비즈니스 로직 레이어
  handler.go      Gin HTTP 핸들러
  grpc_handler.go gRPC 핸들러 (UseCase 재사용)
```

### [go-sample-stdlib-http](./go-sample-stdlib-http)

외부 프레임워크 없이 `net/http`만으로 구현한 REST API 서버입니다.
Gin이 내부에서 처리하는 라우팅, 미들웨어, JSON 인코딩을 직접 작성합니다.

```
cmd/server/   ← 서버 진입점, loggingMiddleware 수동 연결
handler/      Store 인터페이스 + HTTP 핸들러
store/        인메모리 저장소 (sync.RWMutex)
```

### [go-sample-cli](./go-sample-cli)

Cobra를 사용한 CLI 아이템 관리 도구입니다.
`PersistentPreRunE`로 공유 리소스를 초기화하는 Cobra 관용 패턴을 보여줍니다.

```
cmd/          서브커맨드 (list, create, get, delete)
store/        JSON 파일 기반 영속화
```

### [go-sample-tcp](./go-sample-tcp)

`net` 패키지만으로 구현한 TCP 채팅 서버입니다.
클라이언트마다 고루틴을 하나씩 생성하고, `sync.Mutex`로 공유 상태를 보호하는 패턴을 보여줍니다.

```
cmd/server/   hub 구조체, 고루틴 per 연결, broadcast
cmd/client/   송수신 고루틴 분리
```

---

## 실행 방법

각 프로젝트 디렉토리에서 독립적으로 실행합니다.

```bash
# HTTP 서버
cd go-sample-http
go run ./cmd/server

# stdlib HTTP 서버
cd go-sample-stdlib-http
go run ./cmd/server

# CLI
cd go-sample-cli
go run . --help

# TCP 채팅 서버
cd go-sample-tcp
go run ./cmd/server   # 터미널 1
go run ./cmd/client   # 터미널 2, 3 ...
```
