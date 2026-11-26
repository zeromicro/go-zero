# 예제 05: 다중 서비스(`--multiple`)

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 하나의 proto 파일에서 여러 RPC 서비스를 생성하는 방법을 보여줍니다.

## proto 정의

두 proto 파일이 동일한 `go_package`를 공유합니다.

```protobuf
option go_package = "example.com/demo/pb";
```

- `shared.proto` — 공유 메시지 타입(`Meta`)을 정의합니다.
- `multi.proto` — **두 개의** 서비스 `SearchService`와 `NotifyService`를 정의합니다.

proto 파일에 `service` 블록이 둘 이상 포함된 경우 `-m`(또는 `--multiple`) 플래그가 필요합니다.

## 생성 명령

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 `-m` 플래그와 함께 코드를 생성합니다.

```bash
goctl rpc protoc multi.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . \
  -m
```

생성되는 디렉터리 구조:

```
output/
├── client
│   ├── notifyservice
│   │   └── notifyservice.go
│   └── searchservice
│       └── searchservice.go
├── etc
│   └── multisvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── notifyservice
│   │   │   └── notifylogic.go
│   │   └── searchservice
│   │       └── searchlogic.go
│   ├── server
│   │   ├── notifyservice
│   │   │   └── notifyserviceserver.go
│   │   └── searchservice
│   │       └── searchserviceserver.go
│   └── svc
│       └── servicecontext.go
├── multisvc.go
└── pb
    ├── multi.pb.go
    ├── multi_grpc.pb.go
    └── shared.pb.go
```

## 핵심 사항

- `-m`(또는 `--multiple`) 플래그는 다중 서비스 모드를 활성화합니다.
- 다중 모드에서는 `client/`가 서비스별 하위 디렉터리를 포함합니다. `logic/`과 `server/`도 서비스 이름별로 그룹화됩니다.
- 두 서비스는 하나의 진입점(`multisvc.go`)과 설정을 공유합니다.
- `--multiple`이 없으면 goctl은 proto 파일당 하나의 `service` 블록만 허용합니다.
- 모든 서비스는 동일한 `config.go`와 `servicecontext.go`를 공유합니다.
