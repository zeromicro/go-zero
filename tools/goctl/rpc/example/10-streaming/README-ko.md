# 예제 10: 스트리밍 RPC

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 서버 스트리밍, 클라이언트 스트리밍, 양방향 스트리밍이라는 세 가지 gRPC 스트리밍 패턴을 모두 보여줍니다.

## Proto 정의

`stream.proto`는 각 스트리밍 패턴을 보여주는 세 개의 RPC 메서드를 정의합니다.

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
goctl rpc protoc stream.proto \
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
│   └── streamsvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── bidistreamlogic.go
│   │   ├── clientstreamlogic.go
│   │   └── serverstreamlogic.go
│   ├── server
│   │   └── streamserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── stream.pb.go
│   └── stream_grpc.pb.go
├── streamservice
│   └── streamservice.go
└── streamsvc.go
```

## 핵심 사항

- 세 가지 스트리밍 패턴을 지원합니다. 서버 스트리밍(응답에 `stream`), 클라이언트 스트리밍(요청에 `stream`), 양방향 스트리밍(양쪽 모두 `stream`).
- goctl은 각 스트리밍 RPC 메서드마다 별도의 logic 파일을 생성합니다.
- 스트리밍 클라이언트 코드는 자동 생성되지 않으므로 gRPC 클라이언트를 직접 사용하세요.
