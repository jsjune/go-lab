# go-sample-http

Go + Gin으로 만든 CRUD REST API 샘플입니다.
HTTP와 gRPC 두 가지 서버를 제공하며, 둘 다 동일한 비즈니스 로직(UseCase)을 재사용합니다.
Spring Boot 경험자를 위해 익숙한 개념과 비교해서 설명합니다.

---

## 목차

1. [프로젝트 구조](#프로젝트-구조)
2. [Spring vs Go 개념 비교](#spring-vs-go-개념-비교)
3. [설정 (환경변수)](#설정-환경변수)
4. [빌드 & 실행](#빌드--실행)
5. [HTTP API 명세](#http-api-명세)
6. [gRPC](#grpc)
7. [의존성](#의존성)
7. [도메인이 복잡해질 때 구조 진화](#도메인이-복잡해질-때-구조-진화)
8. [Spring 개발자를 위한 팁](#spring-개발자를-위한-팁)

---

## 프로젝트 구조

```
go-sample-http/
├── cmd/
│   ├── server/
│   │   └── main.go           # HTTP 서버 조립 & 실행
│   └── grpc/
│       └── main.go           # gRPC 서버 조립 & 실행
├── proto/
│   └── item.proto            # Protocol Buffer 서비스 정의
├── gen/
│   └── item/                 # protoc가 생성한 코드 (수동 수정 금지)
│       ├── item.pb.go
│       └── item_grpc.pb.go
├── config/
│   └── config.go             # 설정 로드 (Spring의 application.yaml + @ConfigurationProperties)
├── internal/                 # 외부 모듈에서 import 불가 (컴파일러 강제)
│   ├── db/
│   │   └── db.go             # DB 연결 & 테이블 생성 (Spring의 DataSource + schema.sql)
│   └── item/                 # item 도메인
│       ├── entity.go         # 도메인 구조체 (Spring의 @Entity)
│       ├── dto.go            # 요청/응답 DTO (CreateRequest, UpdateRequest)
│       ├── repository.go     # Repository 인터페이스 + SQLite 구현체 (Spring의 JpaRepository)
│       ├── usecase.go        # UseCase 인터페이스 + 비즈니스 로직 (Spring의 @Service)
│       ├── handler.go        # HTTP 핸들러 (Spring의 @RestController)
│       └── grpc_handler.go   # gRPC 핸들러 (HTTP handler와 동일한 UseCase 재사용)
├── .env                      # 로컬 환경변수 (Spring의 application-local.yaml)
├── .env.example              # 환경변수 명세 (팀 공유용)
├── Dockerfile
└── go.mod                    # 의존성 관리 (Maven의 pom.xml / Gradle의 build.gradle)
```

### 레이어 흐름

```
HTTP 요청
    ↓
handler.go      (입력 파싱 & 응답 변환)   ← Spring @RestController
    ↓
usecase.go      (비즈니스 규칙 처리)      ← Spring @Service
    ↓
repository.go   (DB 쿼리)                ← Spring JpaRepository
    ↓
SQLite (app.db)
```

`cmd/server/main.go`와 `cmd/grpc/main.go` 모두 로직 없이 조립만 담당합니다.
UseCase는 HTTP와 gRPC가 **공유**합니다.

```go
// 동일한 UseCase를 HTTP handler와 gRPC handler 둘 다에 주입
repo    := item.NewRepository(db)
uc      := item.NewUseCase(repo)

// HTTP 서버 (cmd/server/main.go)
httpHandler := item.NewHandler(uc)

// gRPC 서버 (cmd/grpc/main.go)
grpcHandler := item.NewGRPCHandler(uc)
```

---

## Spring vs Go 개념 비교

| Spring | Go (이 프로젝트) | 설명 |
|--------|-----------------|------|
| `@SpringBootApplication` + `main()` | `cmd/server/main.go` | 앱 시작점 & 레이어 조립 |
| `pom.xml` / `build.gradle` | `go.mod` | 의존성 관리 |
| `@RestController` | `handler.go` | HTTP 요청 처리 |
| `@Service` | `usecase.go` | 비즈니스 로직 레이어 |
| `JpaRepository` | `repository.go` | DB 접근 레이어 |
| `@Entity` | `entity.go`의 `Item` struct | 도메인 모델 |
| Request DTO | `dto.go`의 `CreateRequest` struct | 요청 바디 |
| `application.yaml` | `.env` + `config/config.go` | 설정 |
| H2 (내장 DB) | SQLite (`app.db`) | 파일 기반 내장 DB |
| Spring MVC | Gin | HTTP 프레임워크 |
| `@Valid` | `binding:"required"` 태그 | 요청 유효성 검사 |
| `@Autowired` / 생성자 주입 | `NewXxx()` 생성자 함수 | 의존성 주입 |
| `ResponseEntity` | `c.JSON(statusCode, body)` | HTTP 응답 |
| `package-private` | `internal/` 디렉토리 | 접근 범위 제한 (컴파일러 강제) |
| `mvn package` → jar | `go build` → 단일 실행 파일 | 배포 산출물 |

> Go는 프레임워크가 자동으로 해주는 것이 적습니다. Spring이 어노테이션으로 처리하는 것들을
> Go에서는 명시적으로 코드로 작성합니다. 대신 동작을 쉽게 추적할 수 있습니다.

---

## 설정 (환경변수)

Spring의 `application.yaml` 대신 환경변수를 사용합니다.

| 환경변수 | 기본값 | 설명 |
|----------|--------|------|
| `APP_PORT` | `8080` | HTTP 서버 포트 |
| `APP_GRPC_PORT` | `9090` | gRPC 서버 포트 |
| `APP_DB_PATH` | `./app.db` | SQLite 파일 경로 |
| `GIN_MODE` | `debug` | `debug` / `release` (Spring Profile과 유사) |

로컬에서는 `.env` 파일에 작성하면 자동으로 읽습니다.

```bash
# .env
APP_PORT=8080
APP_GRPC_PORT=9090
APP_DB_PATH=./app.db
GIN_MODE=debug
```

우선순위: **OS 환경변수 > .env 파일 > 코드 기본값**

---

## 빌드 & 실행

### 개발 중 (빌드 없이 바로 실행)

```bash
# HTTP 서버
go run ./cmd/server

# gRPC 서버 (별도 터미널)
go run ./cmd/grpc
```

Spring의 `./gradlew bootRun`과 동일합니다. 실행 파일은 생성되지 않습니다.

### 실행 파일(exe) 빌드

```bash
# Windows
go build -o server.exe     ./cmd/server   # HTTP 서버
go build -o grpc-server.exe ./cmd/grpc    # gRPC 서버

# Linux / Mac
go build -o server      ./cmd/server
go build -o grpc-server ./cmd/grpc
```

Spring의 `mvn package`로 jar 하나가 나오는 것처럼,
Go는 `go build`로 **단일 실행 파일** 하나가 나옵니다.
JVM 설치 없이 그 파일 하나만 서버에 올리면 바로 실행됩니다.

```bash
# 빌드 후 실행
./server.exe          # Windows
./server              # Linux / Mac
```

### 크로스 컴파일 (다른 OS용 빌드)

Go는 내 머신에서 다른 OS용 바이너리를 빌드할 수 있습니다.

```bash
# Windows 머신에서 Linux 서버용 바이너리 생성
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o server ./cmd/server

# Windows 머신에서 Mac용
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o server ./cmd/server
```

| `GOOS` | `GOARCH` | 대상 |
|--------|----------|------|
| `linux` | `amd64` | 일반 Linux 서버 |
| `linux` | `arm64` | AWS Graviton / Apple M 시리즈 |
| `darwin` | `amd64` | Intel Mac |
| `darwin` | `arm64` | Apple Silicon Mac |
| `windows` | `amd64` | Windows 64bit |

### 핫리로드 (Spring DevTools 유사)

Go는 기본적으로 핫리로드가 없습니다. `air`를 사용하면 파일 변경 시 자동으로 재시작합니다.

```bash
go install github.com/air-verse/air@latest
air
```

### Docker 실행

```bash
# 이미지 빌드
docker build -t go-sample-http .

# 실행 (DB 파일을 로컬에 유지)
docker run -p 8080:8080 \
  -e GIN_MODE=release \
  -v $(pwd)/data:/app \
  go-sample-http
```

---

## HTTP API 명세

Base URL: `http://localhost:8080`

### Health Check

```
GET /health
```

```json
{ "status": "ok", "time": "2026-06-03T02:26:29Z" }
```

### Items CRUD

| Method | URL | 설명 | Status |
|--------|-----|------|--------|
| `GET` | `/items` | 전체 목록 | 200 |
| `GET` | `/items/:id` | 단건 조회 | 200 / 404 |
| `POST` | `/items` | 생성 | 201 |
| `PUT` | `/items/:id` | 수정 | 200 / 404 |
| `DELETE` | `/items/:id` | 삭제 | 204 / 404 |

**요청 바디 (POST / PUT):**

```json
{
  "name": "apple",
  "description": "a red fruit"
}
```

**응답 (단건):**

```json
{
  "id": 1,
  "name": "apple",
  "description": "a red fruit",
  "created_at": "2026-06-03T02:32:16Z"
}
```

---

## gRPC

### 개요

HTTP와 별개로 gRPC 서버를 포트 `9090`에서 제공합니다.
`proto/item.proto`에 서비스를 정의하고, `protoc`로 Go 코드를 자동 생성합니다.

```
proto/item.proto  →  protoc  →  gen/item/item.pb.go
                               gen/item/item_grpc.pb.go
```

HTTP handler와 gRPC handler가 **동일한 UseCase**를 공유하므로
비즈니스 로직 중복 없이 두 프로토콜을 동시에 지원합니다.

```
HTTP  요청 → handler.go      ↘
                              usecase.go → repository.go → DB
gRPC  요청 → grpc_handler.go ↗
```

### proto 코드 재생성

`.proto` 파일을 수정했을 때의 전체 과정(protoc 설치, 코드 생성, 생성된 코드 이해, 필드 번호 주의사항)은 [docs/grpc-codegen.md](docs/grpc-codegen.md)를 참고하세요.

### gRPC 서비스 명세

서비스명: `item.ItemService` / 포트: `9090`

| RPC | 요청 | 응답 | 설명 |
|-----|------|------|------|
| `ListItems` | `{}` | `{items:[...]}` | 전체 목록 |
| `GetItem` | `{id}` | `Item` | 단건 조회 |
| `CreateItem` | `{name, description}` | `Item` | 생성 |
| `UpdateItem` | `{id, name, description}` | `Item` | 수정 |
| `DeleteItem` | `{id}` | `{}` | 삭제 |

### grpcurl로 테스트

grpcurl 설치, 서비스 명세, 전체 CRUD 예시는 [docs/grpc-codegen.md](docs/grpc-codegen.md)를 참고하세요.

---

## 의존성

| 라이브러리 | 역할 | Spring 대응 |
|-----------|------|------------|
| `github.com/gin-gonic/gin` | HTTP 프레임워크 | Spring MVC |
| `modernc.org/sqlite` | 순수 Go SQLite 드라이버 (CGO 불필요) | H2 |
| `github.com/joho/godotenv` | `.env` 파일 로드 | - |
| `google.golang.org/grpc` | gRPC 프레임워크 | Spring gRPC |
| `google.golang.org/protobuf` | Protocol Buffer 런타임 | - |

```bash
# 의존성 추가
go get github.com/some/library   # Maven의 mvn install

# 미사용 의존성 정리
go mod tidy                      # Maven의 dependency:analyze
```

---

## 도메인이 복잡해질 때 구조 진화

### 1단계 — 지금 구조 (단순 CRUD)

```
internal/
└── item/
    ├── entity.go
    ├── dto.go
    ├── repository.go
    ├── usecase.go
    └── handler.go
```

### 2단계 — 도메인 로직이 생길 때

비즈니스 규칙이 복잡해지면 usecase가 비대해집니다.
이때 도메인 객체(entity)에 행위를 부여하고 Value Object를 분리합니다.

```
internal/
└── order/
    ├── entity.go        # Order 구조체 + 도메인 메서드 (order.Cancel() 등)
    ├── vo/              # Value Object (Spring의 @Embedded 타입)
    │   └── money.go
    ├── repository.go
    ├── usecase.go
    └── handler.go
```

```go
// entity.go — 도메인 로직을 struct 메서드로 (Spring의 Rich Domain Model)
func (o *Order) Cancel() error {
    if o.Status == "shipped" {
        return errors.New("shipped 상태는 취소 불가")
    }
    o.Status = "cancelled"
    return nil
}
```

### 3단계 — 도메인 간 의존이 생길 때

`order`가 `user`, `product`를 참조해야 할 때
도메인 간 직접 import를 금지하고 인터페이스로 연결합니다.

```
internal/
├── user/
├── product/
└── order/
    ├── entity.go
    ├── repository.go
    ├── usecase.go
    ├── handler.go
    └── dependency.go    # 다른 도메인에서 필요한 기능만 인터페이스로 선언
```

```go
// order/dependency.go
// order 패키지가 user 패키지를 직접 import하지 않고
// 필요한 기능만 인터페이스로 선언 → 순환 의존 방지
type UserReader interface {
    FindByID(id int64) (*User, error)
}

// order/usecase.go
type useCase struct {
    repo      Repository
    userReader UserReader  // cmd/server/main.go에서 user.Repository를 주입
}
```

Spring의 서비스 간 인터페이스 분리 또는 `@FeignClient`와 유사한 개념입니다.

### 4단계 — 횡단 관심사가 생길 때

인증, 로깅처럼 여러 도메인에 걸친 관심사는 `shared/`로 분리합니다.

```
internal/
├── user/
├── order/
├── product/
└── shared/                # Spring의 공통 컴포넌트
    ├── middleware/         # Spring의 Filter / Interceptor
    │   ├── auth.go
    │   └── logger.go
    ├── apperror/           # 공통 에러 타입
    └── pagination/         # 공통 페이지네이션 구조체
```

### 최종 형태 — 대규모 서비스

```
my-go-app/
├── cmd/
│   ├── server/main.go     # HTTP API 서버
│   └── worker/main.go     # 배치 / 이벤트 소비자 (Spring Batch / @Scheduled)
├── config/
├── internal/
│   ├── user/
│   ├── order/
│   │   ├── entity.go
│   │   ├── vo/
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   ├── handler.go
│   │   └── dependency.go
│   ├── product/
│   └── shared/
│       ├── middleware/
│       ├── apperror/
│       └── pagination/
└── pkg/                   # 외부 모듈에서도 쓸 수 있는 범용 유틸
    └── timeutil/
```

> **Go다운 방식:** 지금 구조에서 **필요할 때만** 단계적으로 추가합니다.
> Spring처럼 처음부터 모든 레이어를 만들지 않습니다.

---

## Spring 개발자를 위한 팁

**어노테이션이 없습니다.**
Spring의 `@Autowired`, `@Transactional`, `@Valid` 같은 어노테이션 대신
모든 것을 명시적인 함수 호출로 작성합니다.

**DI 컨테이너가 없습니다.**
Spring IoC 컨테이너 대신 `cmd/server/main.go`에서 직접 객체를 생성하고 연결합니다.

```go
// Spring: @Autowired가 자동으로 주입
// Go: main.go에서 직접 연결
repo    := item.NewRepository(db)
uc      := item.NewUseCase(repo)
handler := item.NewHandler(uc)
```

**에러는 반환값으로 처리합니다.**
Spring의 `try-catch` / `@ExceptionHandler` 대신 Go는 에러를 반환값으로 다룹니다.

```go
item, err := repo.FindByID(id)
if err != nil {
    // 에러 처리
}
```

**인터페이스는 명시적으로 구현하지 않습니다.**
Spring의 `implements` 없이 메서드 시그니처가 맞으면 자동으로 인터페이스를 구현한 것으로 봅니다.

```go
// Repository 인터페이스 선언
type Repository interface {
    FindAll() ([]Item, error)
}

// implements 키워드 없이 메서드만 맞으면 자동으로 구현체로 인정
type sqliteRepository struct{ db *sql.DB }
func (r *sqliteRepository) FindAll() ([]Item, error) { ... }
```

**배포 산출물이 단순합니다.**
Spring은 JVM + jar가 필요하지만 Go는 단일 실행 파일 하나로 배포합니다.

```
Spring:  서버에 JDK 설치 → jar 업로드 → java -jar app.jar
Go:      server 파일 하나 업로드 → ./server
```
