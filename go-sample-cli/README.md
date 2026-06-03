# go-sample-cli

[Cobra](https://github.com/spf13/cobra)를 사용한 CLI 아이템 관리 도구 샘플입니다.
데이터는 `~/.item-cli/data.json`에 JSON 파일로 영속화됩니다.

---

## 폴더 구조

```
go-sample-cli/
├── cmd/
│   ├── root.go     # 루트 커맨드 + PersistentPreRunE로 Store 초기화, 서브커맨드 등록
│   ├── list.go     # list 서브커맨드
│   ├── create.go   # create 서브커맨드 (플래그 바인딩 예시)
│   ├── get.go      # get 서브커맨드 (위치 인자 예시)
│   └── delete.go   # delete 서브커맨드
├── store/
│   └── store.go    # JSON 파일 읽기/쓰기, 슬라이스 CRUD
├── main.go         # cmd.Execute() 위임
└── go.mod
```

---

## 실행

```bash
# 빌드 없이 바로 실행
go run . list
go run . create --name "apple" --description "red fruit"
go run . create -n "banana"
go run . get 1
go run . delete 1

# 도움말
go run . --help
go run . create --help
```

## 빌드 후 실행

```bash
go build -o item .

./item list
./item create -n "apple" -d "red fruit"
./item get 1
./item delete 1
```

---

## 학습 포인트

| 파일 | 내용 |
|------|------|
| `main.go` | Cobra 관용 진입점 — 로직 없이 `cmd.Execute()` 한 줄만 |
| `cmd/root.go` | `PersistentPreRunE`로 공유 리소스(Store)를 한 번만 초기화 |
| `cmd/create.go` | `Flags().StringVarP()`로 플래그 바인딩, `RunE`로 에러 반환 |
| `cmd/get.go` | `Args: cobra.ExactArgs(1)`로 위치 인자 개수 강제 |
| `store/store.go` | JSON 파일 읽기/쓰기, `os.UserHomeDir()`로 OS 독립 경로 처리 |

---

## 아키텍처

```
main.go
  └── cmd.Execute()
        └── rootCmd
              PersistentPreRunE: store.New() → globalStore
              ├── listCmd   → globalStore.List()
              ├── createCmd → globalStore.Create()
              ├── getCmd    → globalStore.Get()
              └── deleteCmd → globalStore.Delete()
```

`PersistentPreRunE`는 모든 서브커맨드 실행 전에 먼저 호출됩니다.
Store 초기화(파일 읽기)를 각 커맨드가 아닌 루트에서 한 번만 수행합니다.

---

## 데이터 저장 위치

| OS | 경로 |
|----|------|
| Windows | `C:\Users\<사용자>\.item-cli\data.json` |
| macOS / Linux | `~/.item-cli/data.json` |

---

## 의존성

| 라이브러리 | 역할 |
|-----------|------|
| `github.com/spf13/cobra` | CLI 프레임워크 — 서브커맨드, 플래그, 도움말 자동 생성 |
