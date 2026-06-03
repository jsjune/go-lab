# gRPC 코드 자동 생성 가이드

proto 파일을 작성하고 Go 코드를 자동 생성하는 전체 과정을 설명합니다.

---

## 전체 흐름

```
1. proto 파일 작성       proto/item.proto
         ↓ protoc 실행
2. Go 코드 자동 생성     gen/item/item.pb.go          ← 메시지 struct
                         gen/item/item_grpc.pb.go     ← 서버/클라이언트 인터페이스
         ↓ 인터페이스 확인 후
3. 구현체 작성           internal/item/grpc_handler.go
```

---

## 1단계 — 사전 설치

### protoc (컴파일러 본체)

**Windows:**

```powershell
# winget으로 설치
winget install Google.Protobuf

# 또는 직접 다운로드
# https://github.com/protocolbuffers/protobuf/releases
# protoc-xx.x-win64.zip 다운로드 후 bin/protoc.exe를 PATH에 추가
```

**Mac:**

```bash
brew install protobuf
```

**Linux:**

```bash
apt install -y protobuf-compiler
```

설치 확인:

```bash
protoc --version
# libprotoc 29.x
```

### Go 플러그인 (Go 코드 생성용)

```bash
# 메시지 struct 생성 플러그인
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# gRPC 인터페이스 생성 플러그인
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

설치 확인:

```bash
protoc-gen-go --version
protoc-gen-go-grpc --version
```

> 플러그인은 Go bin 경로(`~/go/bin`)에 설치됩니다.
> 이 경로가 PATH에 없으면 protoc 실행 시 플러그인을 찾지 못합니다.
>
> ```bash
> # PATH 확인 및 추가 (bash/zsh)
> export PATH="$PATH:$(go env GOPATH)/bin"
>
> # PowerShell
> $env:PATH += ";$(go env GOPATH)\bin"
> ```

---

## 2단계 — proto 파일 작성

`proto/item.proto`

```protobuf
syntax = "proto3";

package item;

// 생성될 Go 코드의 import 경로
option go_package = "go-sample-http/gen/item";

// 서비스 정의 (어떤 RPC 메서드를 제공할지)
service ItemService {
  rpc ListItems   (ListItemsRequest)   returns (ListItemsResponse);
  rpc GetItem     (GetItemRequest)     returns (ItemResponse);
  rpc CreateItem  (CreateItemRequest)  returns (ItemResponse);
  rpc UpdateItem  (UpdateItemRequest)  returns (ItemResponse);
  rpc DeleteItem  (DeleteItemRequest)  returns (DeleteItemResponse);
}

// 메시지 정의 (요청/응답 타입)
message ListItemsRequest {}

message ListItemsResponse {
  repeated ItemResponse items = 1;   // repeated = 배열 (Java의 List)
}

message GetItemRequest {
  int64 id = 1;
}

message CreateItemRequest {
  string name        = 1;
  string description = 2;
}

message UpdateItemRequest {
  int64  id          = 1;
  string name        = 2;
  string description = 3;
}

message DeleteItemRequest {
  int64 id = 1;
}

message DeleteItemResponse {}

message ItemResponse {
  int64  id          = 1;
  string name        = 2;
  string description = 3;
  string created_at  = 4;
}
```

### proto 파일 문법 핵심

```protobuf
message MyMessage {
  string  name  = 1;   // = 1, = 2 는 필드 번호 (순서 ID, 이름 변경해도 유지됨)
  int64   id    = 2;
  bool    active = 3;

  repeated string tags = 4;   // 배열
  optional string memo = 5;   // nullable (proto3에서는 기본이 optional)
}
```

| proto 타입 | Go 타입 | Java 타입 |
|-----------|---------|-----------|
| `string` | `string` | `String` |
| `int32` | `int32` | `int` |
| `int64` | `int64` | `long` |
| `bool` | `bool` | `boolean` |
| `repeated X` | `[]X` | `List<X>` |

---

## 3단계 — 코드 생성 실행

프로젝트 루트에서 실행합니다.

```bash
protoc \
  --go_out=gen/item \
  --go_opt=paths=source_relative \
  --go-grpc_out=gen/item \
  --go-grpc_opt=paths=source_relative \
  --proto_path=proto \
  proto/item.proto
```

### 옵션 설명

| 옵션 | 의미 |
|------|------|
| `--proto_path=proto` | proto 파일을 찾을 디렉토리 |
| `--go_out=gen/item` | `item.pb.go` 출력 위치 |
| `--go_opt=paths=source_relative` | 출력 경로를 proto 파일 기준으로 |
| `--go-grpc_out=gen/item` | `item_grpc.pb.go` 출력 위치 |
| `--go-grpc_opt=paths=source_relative` | 위와 동일 |
| `proto/item.proto` | 변환할 파일 |

### 생성 결과

```
gen/item/
├── item.pb.go       ← 메시지 struct (ListItemsRequest, ItemResponse 등)
└── item_grpc.pb.go  ← 서버/클라이언트 인터페이스 (ItemServiceServer 등)
```

---

## 4단계 — 생성된 코드 이해

### item.pb.go — 메시지 struct

proto의 `message`가 Go `struct`로 변환됩니다.

```go
// proto 정의
message ItemResponse {
  int64  id   = 1;
  string name = 2;
}

// ↓ 자동 생성된 Go 코드
type ItemResponse struct {
    Id   int64  `protobuf:"varint,1,..."`
    Name string `protobuf:"bytes,2,..."`
    // ...내부 필드 생략
}
```

### item_grpc.pb.go — 서버 인터페이스

proto의 `service`가 Go `interface`로 변환됩니다.
**개발자는 이 인터페이스를 보고 구현체를 작성합니다.**

```go
// 자동 생성된 서버 인터페이스
type ItemServiceServer interface {
    ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error)
    GetItem(context.Context, *GetItemRequest) (*ItemResponse, error)
    CreateItem(context.Context, *CreateItemRequest) (*ItemResponse, error)
    UpdateItem(context.Context, *UpdateItemRequest) (*ItemResponse, error)
    DeleteItem(context.Context, *DeleteItemRequest) (*DeleteItemResponse, error)
    mustEmbedUnimplementedItemServiceServer()
}

// 자동 생성된 클라이언트 인터페이스 (gRPC 클라이언트 개발 시 사용)
type ItemServiceClient interface {
    ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error)
    // ...
}
```

---

## 5단계 — 구현체 작성

자동 생성된 인터페이스를 구현합니다.

```go
// internal/item/grpc_handler.go (개발자 작성)
type GRPCHandler struct {
    pb.UnimplementedItemServiceServer  // 미구현 메서드 기본값 제공
    uc UseCase
}

func (h *GRPCHandler) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.ItemResponse, error) {
    item, err := h.uc.Get(req.Id)
    if err != nil {
        return nil, status.Error(codes.NotFound, "item not found")
    }
    return &pb.ItemResponse{
        Id:   item.ID,
        Name: item.Name,
    }, nil
}
```

> `UnimplementedItemServiceServer`를 embed하면 구현하지 않은 메서드는
> 자동으로 `codes.Unimplemented` 에러를 반환합니다.
> Spring의 `default` 메서드가 있는 인터페이스와 비슷합니다.

---

## proto 파일 수정 시

서비스나 메시지를 변경할 때마다 3단계 명령어를 다시 실행합니다.

```
proto 수정 → protoc 재실행 → gen/ 코드 갱신 → 구현체 수정
```

### 주의사항 — 필드 번호

```protobuf
message Item {
  int64  id   = 1;
  string name = 2;   // ← 이 번호를 절대 바꾸지 말것
}
```

필드 번호(`= 1`, `= 2`)는 직렬화 시 사용되는 식별자입니다.
이미 배포된 서비스에서 번호를 바꾸면 기존 클라이언트가 데이터를 읽지 못합니다.
필드를 제거할 때는 번호를 `reserved`로 예약합니다.

```protobuf
message Item {
  reserved 2;          // 이 번호는 다시 사용 불가
  reserved "name";     // 이 이름도 다시 사용 불가
  int64 id = 1;
}
```

---

## 서비스 명세

서비스명: `item.ItemService` / 포트: `9090`

| RPC | 요청 | 응답 | 설명 |
|-----|------|------|------|
| `ListItems` | `{}` | `{items:[...]}` | 전체 목록 |
| `GetItem` | `{id}` | `Item` | 단건 조회 |
| `CreateItem` | `{name, description}` | `Item` | 생성 |
| `UpdateItem` | `{id, name, description}` | `Item` | 수정 |
| `DeleteItem` | `{id}` | `{}` | 삭제 |

---

## 자주 쓰는 명령어 요약

```bash
# proto → Go 코드 생성
protoc \
  --go_out=gen/item --go_opt=paths=source_relative \
  --go-grpc_out=gen/item --go-grpc_opt=paths=source_relative \
  --proto_path=proto proto/item.proto

# gRPC 서버 실행
go run ./cmd/grpc

# grpcurl 설치 (최초 1회)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 서비스 목록 조회 (reflection 활성화 덕분에 proto 파일 없이 가능)
grpcurl -plaintext localhost:9090 list

# CRUD 전체 예시
grpcurl -plaintext localhost:9090 item.ItemService/ListItems

grpcurl -plaintext \
  -d '{"name":"apple","description":"a red fruit"}' \
  localhost:9090 item.ItemService/CreateItem

grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:9090 item.ItemService/GetItem

grpcurl -plaintext \
  -d '{"id":1,"name":"apple","description":"updated"}' \
  localhost:9090 item.ItemService/UpdateItem

grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:9090 item.ItemService/DeleteItem
```
