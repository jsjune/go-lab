# go-sample-stdlib-http

외부 프레임워크 없이 Go 표준 라이브러리(`net/http`)만으로 구현한 REST API 샘플입니다.
`go-sample-http`(Gin 버전)와 동일한 기능을 제공하며, Gin이 내부에서 처리하는 것을 직접 구현합니다.

---

## 폴더 구조

```
go-sample-stdlib-http/
├── cmd/
│   └── server/
│       └── main.go     # 서버 조립 & 실행, loggingMiddleware 적용
├── handler/
│   └── item.go         # Store 인터페이스 정의, HTTP 핸들러
├── store/
│   └── store.go        # 인메모리 저장소 (sync.RWMutex)
└── go.mod
```

---

## 실행

```bash
go run ./cmd/server
```

---

## 학습 포인트

| 파일 | Gin 대응 | 내용 |
|------|----------|------|
| `cmd/server/main.go` | `gin.Default()` + `r.Use()` | `http.NewServeMux()` + 미들웨어를 핸들러로 감싸기 |
| `handler/item.go` | `r.Group("/items")` | `HandleFunc("/items", ...)` + `switch r.Method` 분기 |
| `handler/item.go` | `c.ShouldBindJSON` | `json.NewDecoder(r.Body).Decode` |
| `handler/item.go` | `c.JSON(200, obj)` | `w.Header().Set(...)` + `json.NewEncoder(w).Encode` |
| `store/store.go` | - | `sync.RWMutex` — 읽기 동시 허용, 쓰기 단독 보호 |
| `handler/item.go` | 구체 타입 의존 없음 | `Store` 인터페이스 — 테스트 시 목 주입 가능 |

---

## Gin과 표준 라이브러리 비교

| 기능 | Gin | 표준 라이브러리 |
|------|-----|----------------|
| 라우팅 | `r.GET("/items/:id", h)` | `mux.HandleFunc("/items/", h)` + 경로 수동 파싱 |
| 메서드 분기 | 자동 | `switch r.Method { ... }` |
| JSON 바인딩 | `c.ShouldBindJSON(&req)` | `json.NewDecoder(r.Body).Decode(&req)` |
| JSON 응답 | `c.JSON(200, obj)` | `w.Header().Set(...)` + `json.NewEncoder(w).Encode(obj)` |
| 미들웨어 | `r.Use(fn)` | `loggingMiddleware(mux)` — 핸들러를 감싸는 함수 |
| 경로 파라미터 | `c.Param("id")` | `strings.TrimPrefix(r.URL.Path, "/items/")` |

---

## API

| Method | URL | 설명 | 상태코드 |
|--------|-----|------|---------|
| `GET` | `/health` | 헬스 체크 | 200 |
| `GET` | `/items` | 전체 목록 | 200 |
| `GET` | `/items/:id` | 단건 조회 | 200 / 404 |
| `POST` | `/items` | 생성 | 201 |
| `PUT` | `/items/:id` | 수정 | 200 / 404 |
| `DELETE` | `/items/:id` | 삭제 | 204 / 404 |

**요청 바디 (POST / PUT):**

```json
{ "name": "apple", "description": "red fruit" }
```

---

## 인터페이스 패턴

`handler.ItemHandler`는 `*store.InMemoryStore`에 직접 의존하지 않고 `Store` 인터페이스에 의존합니다.

```go
type Store interface {
    List() []store.Item
    Get(id int64) (*store.Item, error)
    Create(name, description string) store.Item
    Update(id int64, name, description string) (*store.Item, error)
    Delete(id int64) error
}
```

`*store.InMemoryStore`는 이 인터페이스를 자동으로 충족합니다(Go의 암묵적 인터페이스 구현).
테스트 시 목(mock) 구현체를 주입할 수 있습니다.
