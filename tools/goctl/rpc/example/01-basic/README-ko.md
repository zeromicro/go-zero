# 예제 01: 기본 RPC 서비스

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 goctl로 RPC 서비스를 생성하는 가장 간단한 예제입니다.

## proto 정의

외부 import 없이 하나의 서비스와 하나의 RPC 메서드를 가진 `greeter.proto` 파일 하나를 사용합니다.

`go_package`는 전체 모듈 경로를 사용합니다.

```protobuf
option go_package = "example.com/demo/greeter";
```

## 생성 명령

### 방법 1: `goctl rpc new`로 빠르게 시작

```bash
# 한 번의 명령으로 완전한 RPC 프로젝트 생성
goctl rpc new greeter
```

이 명령은 proto 파일과 서비스 코드를 함께 생성합니다.

```
greeter/
├── etc
│   └── greeter.yaml
├── greeter
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeter.proto
├── greeterclient
│   └── greeter.go
└── internal
    ├── config
    │   └── config.go
    ├── logic
    │   └── pinglogic.go
    ├── server
    │   └── greeterserver.go
    └── svc
        └── servicecontext.go
```

### 방법 2: 기존 proto에서 생성

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 코드를 생성합니다.

```bash
goctl rpc protoc greeter.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I .
```

생성되는 디렉터리 구조:

```
output/
├── etc
│   └── greeter.yaml
├── go.mod
├── greeter
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeterclient
│   └── greeter.go
└── internal
    ├── config
    │   └── config.go
    ├── logic
    │   └── sayhellologic.go
    ├── server
    │   └── greeterserver.go
    └── svc
        └── servicecontext.go
```

## 핵심 사항

- 가장 단순한 시나리오입니다. proto 파일 하나, 서비스 하나, RPC 메서드 하나를 사용합니다.
- `go_package`는 상대 경로가 아닌 전체 모듈 경로(`example.com/demo/greeter`)를 사용합니다.
- `--module` 플래그는 goctl에 Go 모듈 이름을 알려줍니다. `--go_opt=module=...`과 `--go-grpc_opt=module=...`은 protoc에 출력 경로에서 모듈 접두사를 제거하라고 알려줍니다.
- `--zrpc_out` 플래그는 goctl이 생성하는 서비스 코드의 출력 위치를 지정합니다.
- `--go_out`과 `--go-grpc_out` 플래그는 protoc가 생성하는 코드의 출력 위치를 지정합니다.
- 비즈니스 로직을 구현하려면 logic 파일(`internal/logic/sayhellologic.go`)을 수정하세요.
