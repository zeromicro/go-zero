# goctl rpc — RPC 코드 생성

[English](README.md) | [中文](README-cn.md) | 한국어

goctl rpc는 `goctl` 스캐폴드의 RPC 서비스 코드 생성 모듈입니다. `.proto` 파일에서 완전한 zRPC 서비스를 생성합니다. proto 정의와 비즈니스 로직만 작성하면 나머지 보일러플레이트 코드는 모두 자동으로 생성됩니다.

## 기능

- **protoc 호환**: protoc와 완전히 호환되며 모든 protoc 인수를 그대로 전달합니다.
- **외부 proto import**: 디렉터리와 패키지를 넘나드는 proto import를 지원하고 전이 의존성을 자동으로 해석합니다.
- **다중 서비스**: 하나의 proto 파일에 여러 서비스를 정의하고, 서비스 이름별로 자동 그룹화합니다.
- **스트리밍 지원**: 서버 스트리밍, 클라이언트 스트리밍, 양방향 스트리밍을 지원합니다.
- **Google well-known types**: `google.protobuf.*` 타입을 자동으로 인식하고 올바른 Go import를 생성합니다.
- **클라이언트 생성**: 자동 생성된 RPC 클라이언트 래퍼 코드를 제공합니다.

## 사전 요구 사항

```bash
# protoc 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 빠른 시작

### 방법 1: 즉시 서비스 만들기

```bash
goctl rpc new greeter
```

완전한 프로젝트 구조가 생성됩니다.

```
greeter/
├── etc/
│   └── greeter.yaml
├── greeter/
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeter.proto
├── greeterclient/
│   └── greeter.go
└── internal/
    ├── config/
    │   └── config.go
    ├── logic/
    │   └── pinglogic.go
    ├── server/
    │   └── greeterserver.go
    └── svc/
        └── servicecontext.go
```

### 방법 2: proto 파일에서 생성하기

1. proto 템플릿을 생성합니다.

```bash
goctl rpc template -o=user.proto
```

2. 출력 디렉터리를 초기화하고 서비스 코드를 생성합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
goctl rpc protoc user.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

---

## 명령 참조

### `goctl rpc protoc`

`.proto` 파일에서 zRPC 서비스 코드를 생성합니다.

```bash
goctl rpc protoc <proto_file> [flags]
```

**예시:**

```bash
# 기본 사용법
goctl rpc protoc user.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .

# 다중 서비스 모드
goctl rpc protoc multi.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I . -m

# 외부 proto import
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I . -I ./shared_protos

# Google well-known types 사용
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

**플래그:**

| 플래그 | 축약 | 타입 | 기본값 | 설명 |
|------|-------|------|---------|-------------|
| `--zrpc_out` | | string | **필수** | zRPC 서비스 코드 출력 디렉터리 |
| `--go_out` | | string | **필수** | protoc Go 코드 출력 디렉터리 |
| `--go-grpc_out` | | string | **필수** | protoc gRPC 코드 출력 디렉터리 |
| `--go_opt` | | string | | protoc-gen-go 옵션(예: `module=example.com/demo`) |
| `--go-grpc_opt` | | string | | protoc-gen-go-grpc 옵션(예: `module=example.com/demo`) |
| `--proto_path` | `-I` | string[] | | proto import 검색 디렉터리(반복 지정 가능) |
| `--multiple` | `-m` | bool | `false` | 다중 서비스 모드 |
| `--client` | `-c` | bool | `true` | RPC 클라이언트 코드 생성 |
| `--style` | | string | `gozero` | 파일 이름 스타일 |
| `--module` | | string | | 사용자 지정 Go module 이름 |
| `--name-from-filename` | | bool | `false` | 서비스 이름에 `package` 이름 대신 파일 이름 사용 |
| `--verbose` | `-v` | bool | `false` | 상세 로그 활성화 |
| `--home` | | string | | goctl 템플릿 디렉터리 |
| `--remote` | | string | | 원격 템플릿 Git 저장소 URL |
| `--branch` | | string | | 원격 템플릿 브랜치 |

### `goctl rpc new`

완전한 RPC 서비스 프로젝트를 빠르게 생성합니다.

```bash
goctl rpc new <service_name> [flags]
```

**플래그:**

| 플래그 | 축약 | 타입 | 기본값 | 설명 |
|------|-------|------|---------|-------------|
| `--style` | | string | `gozero` | 파일 이름 스타일 |
| `--client` | `-c` | bool | `true` | RPC 클라이언트 코드 생성 |
| `--module` | | string | | 사용자 지정 Go module 이름 |
| `--verbose` | `-v` | bool | `false` | 상세 로그 활성화 |
| `--idea` | | bool | `false` | IDE 프로젝트 마커 생성 |
| `--name-from-filename` | | bool | `false` | 서비스 이름에 `package` 이름 대신 파일 이름 사용 |
| `--home` | | string | | goctl 템플릿 디렉터리 |
| `--remote` | | string | | 원격 템플릿 Git 저장소 URL |
| `--branch` | | string | | 원격 템플릿 브랜치 |

### `goctl rpc template`

proto 파일 템플릿을 생성합니다.

```bash
goctl rpc template -o=<output_file> [flags]
```

**플래그:**

| 플래그 | 타입 | 설명 |
|------|------|-------------|
| `-o` | string | 출력 파일 경로(필수) |
| `--home` | string | goctl 템플릿 디렉터리 |
| `--remote` | string | 원격 템플릿 Git 저장소 URL |
| `--branch` | string | 원격 템플릿 브랜치 |

---

## 기능 상세

### 다중 서비스 모드(`--multiple`)

proto 파일에 여러 `service` 정의가 포함된 경우 `--multiple` 플래그가 필요합니다.

```protobuf
service SearchService {
  rpc Search(SearchReq) returns (SearchReply);
}

service NotifyService {
  rpc Notify(NotifyReq) returns (NotifyReply);
}
```

**`--multiple` 사용 시 디렉터리 차이:**

| 기능 | 기본 모드 | `--multiple` 모드 |
|---------|-------------|-------------------|
| proto당 서비스 수 | 정확히 1개 | 1개 이상 |
| 클라이언트 디렉터리 | 서비스 이름 기반 | 고정 `client/` 디렉터리 |
| 코드 구성 | 평면 구조 | 서비스 이름별 그룹화 |

**`--multiple=false`(기본값) 디렉터리 구조:**

```
output/
├── greeterclient/
│   └── greeter.go
├── internal/
│   ├── logic/
│   │   └── sayhellologic.go
│   └── server/
│       └── greeterserver.go
└── ...
```

**`--multiple=true` 디렉터리 구조:**

```
output/
├── client/
│   ├── searchservice/
│   │   └── searchservice.go
│   └── notifyservice/
│       └── notifyservice.go
├── internal/
│   ├── logic/
│   │   ├── searchservice/
│   │   │   └── searchlogic.go
│   │   └── notifyservice/
│   │       └── notifylogic.go
│   └── server/
│       ├── searchservice/
│       │   └── searchserviceserver.go
│       └── notifyservice/
│           └── notifyserviceserver.go
└── ...
```

### 외부 proto import(`--proto_path`)

`-I` / `--proto_path`로 추가 proto 검색 디렉터리를 지정합니다. 지원되는 시나리오는 다음과 같습니다.

- **같은 디렉터리 import**: `import "types.proto";`
- **하위 디렉터리 import**: `import "common/types.proto";`
- **외부 디렉터리 import**: 프로젝트 외부에 있는 proto 파일
- **전이 import**: A가 B를 import하고 B가 C를 import하는 경우 — goctl이 재귀적으로 해석합니다.
- **교차 패키지 import**: 서로 다른 `go_package` 값을 가진 파일에 대해 올바른 Go import를 자동으로 생성합니다.

```bash
# 여러 디렉터리에서 proto 파일 검색
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . -I ./shared_protos -I /path/to/external_protos
```

### 서비스 이름 지정

기본적으로 서비스 이름은 proto의 **`package` 이름**에서 파생됩니다(예: `package user;` → 서비스 이름 `user`). 이를 통해 여러 proto 파일이 같은 `package`를 공유할 수 있습니다.

```
protos/
├── user_base.proto      # package user;
├── user_auth.proto      # package user;
└── user_profile.proto   # package user;
```

세 파일은 모두 하나의 `user` 서비스로 생성됩니다.

서비스 이름에 proto 파일 이름을 사용하려면(레거시 동작) `--name-from-filename` 플래그를 추가하세요.

### 스트리밍 RPC

세 가지 gRPC 스트리밍 패턴을 모두 지원합니다.

```protobuf
service StreamService {
  rpc ServerStream(Req) returns (stream Reply);       // 서버 스트리밍
  rpc ClientStream(stream Req) returns (Reply);       // 클라이언트 스트리밍
  rpc BidiStream(stream Req) returns (stream Reply);  // 양방향 스트리밍
}
```

### Google Well-Known Types

goctl은 Google protobuf well-known types를 자동으로 인식하고 처리합니다.

| Proto Type | Go Type |
|-----------|---------|
| `google.protobuf.Empty` | `emptypb.Empty` |
| `google.protobuf.Timestamp` | `timestamppb.Timestamp` |
| `google.protobuf.Duration` | `durationpb.Duration` |
| `google.protobuf.Any` | `anypb.Any` |
| `google.protobuf.Struct` | `structpb.Struct` |
| `google.protobuf.FieldMask` | `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` | `wrapperspb.*Value` |

이 타입들은 RPC 파라미터 타입으로 직접 사용할 수 있으며, goctl이 올바른 import를 자동으로 생성합니다.

---

## 예제

모든 생성 시나리오를 다루는 10개의 완전한 예제는 [example/](example/) 디렉터리를 참고하세요.

| # | 예제 | 시나리오 |
|---|---------|----------|
| 01 | [기본 서비스](example/01-basic/) | 단일 서비스, import 없음 |
| 02 | [동일 디렉터리 import](example/02-import-sibling/) | 같은 디렉터리에서 import |
| 03 | [하위 디렉터리 import](example/03-import-subdir/) | 하위 디렉터리에서 import |
| 04 | [전이 import](example/04-transitive-import/) | A → B → C 의존성 체인 |
| 05 | [다중 서비스](example/05-multiple-services/) | `--multiple` 모드 |
| 06 | [Well-known 타입](example/06-wellknown-types/) | 메시지에서 Timestamp 등 사용 |
| 07 | [외부 proto(동일 패키지)](example/07-external-proto-same-pkg/) | 외부 proto, 같은 go_package |
| 08 | [외부 proto(다른 패키지)](example/08-external-proto-diff-pkg/) | 외부 proto, 다른 go_package |
| 09 | [Google well-known types를 파라미터로 사용](example/09-google-types-as-rpc/) | Empty/Timestamp를 RPC 파라미터로 사용 |
| 10 | [스트리밍](example/10-streaming/) | 서버/클라이언트/양방향 스트리밍 |
