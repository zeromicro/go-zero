# 예제 09: RPC 파라미터로 Google Types 사용

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 Google protobuf well-known types를 메시지 필드뿐 아니라 RPC 요청 또는 응답 타입으로 **직접** 사용하는 방법을 보여줍니다.

## Proto 정의

`service.proto`는 `google.protobuf.Empty`와 `google.protobuf.Timestamp`를 RPC 요청/응답 타입으로 직접 사용합니다.

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
goctl rpc protoc service.proto \
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
│   └── healthsvc.yaml
├── go.mod
├── healthservice
│   └── healthservice.go
├── healthsvc.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── gettimelogic.go
│   │   └── pinglogic.go
│   ├── server
│   │   └── healthserviceserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    ├── service.pb.go
    └── service_grpc.pb.go
```

## 핵심 사항

- Google well-known types(`google.protobuf.Empty`, `google.protobuf.Timestamp`)를 메시지 필드뿐 아니라 RPC 요청/응답 타입으로 직접 사용합니다.
- goctl은 이를 Go 타입(`emptypb.Empty`, `timestamppb.Timestamp`)으로 올바르게 매핑하고 적절한 import를 생성합니다.
- well-known types를 메시지 필드로 사용하는 예제 06과는 다릅니다.
