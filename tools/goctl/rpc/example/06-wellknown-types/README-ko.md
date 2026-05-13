# 예제 06: Well-Known Types

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 Google protobuf well-known types(`Timestamp`, `Duration`, `Any`)를 메시지 필드로 사용하는 방법을 보여줍니다.

## Proto 정의

`events.proto`는 `google.protobuf.Timestamp`를 메시지 필드 타입으로 사용합니다.

`go_package`는 전체 모듈 경로를 사용합니다.

```protobuf
option go_package = "example.com/demo/pb";
```

## 생성 명령

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 코드를 생성합니다.

```bash
goctl rpc protoc events.proto \
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
│   └── eventsvc.yaml
├── eventservice
│   └── eventservice.go
├── eventsvc.go
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── createeventlogic.go
│   │   └── listeventslogic.go
│   ├── server
│   │   └── eventserviceserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    ├── events.pb.go
    └── events_grpc.pb.go
```

## 핵심 사항

- Google well-known types(`google.protobuf.Timestamp`, `google.protobuf.Duration`, `google.protobuf.Any`)를 메시지 필드로 사용합니다.
- goctl은 well-known types를 Go import(`timestamppb`, `durationpb`, `anypb` 등)로 자동 매핑합니다.
- protoc가 올바르게 설치되어 있다면 well-known types에는 추가 `--proto_path`가 필요하지 않습니다.
