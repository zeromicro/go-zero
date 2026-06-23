# 예제 02: 같은 디렉터리의 proto 파일 import

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 같은 디렉터리에 있는 proto 파일을 import하는 방법을 보여줍니다.

## proto 정의

같은 디렉터리의 두 proto 파일이 동일한 `go_package`를 공유합니다.

- `types.proto` — 공유 메시지 타입(`User`)을 정의합니다.
- `user.proto` — RPC 서비스를 정의하고 `types.proto`를 import합니다.

두 파일은 전체 모듈 경로를 가진 동일한 `go_package`를 사용합니다.

```protobuf
option go_package = "example.com/demo/pb";
```

`user.proto`는 다음과 같이 `types.proto`를 import합니다.

```protobuf
import "types.proto";
```

## 생성 명령

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 코드를 생성합니다.

```bash
goctl rpc protoc user.proto \
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
│   └── usersvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── createuserlogic.go
│   │   └── getuserlogic.go
│   ├── server
│   │   └── userserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── types.pb.go
│   ├── user.pb.go
│   └── user_grpc.pb.go
├── userservice
│   └── userservice.go
└── usersvc.go
```

## 핵심 사항

- 두 proto 파일(`user.proto`와 `types.proto`)은 동일한 `go_package = "example.com/demo/pb"`를 공유하며 하나의 Go 패키지로 컴파일됩니다.
- `user.proto`는 `import "types.proto"`로 `types.proto`를 import합니다.
- 여러 proto 파일이 동일한 `go_package`를 공유하면 하나의 Go 패키지로 컴파일됩니다.
- `service` 정의를 포함한 proto 파일만 `goctl rpc protoc`에 전달하면 됩니다.
- import된 proto는 protoc가 자동으로 컴파일하고 goctl이 해결합니다.
