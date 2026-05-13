# 예제 04: 전이 import

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 A가 B를 import하고 B가 C를 import하는 전이 proto import를 보여줍니다.

## proto 정의

세 proto 파일이 전이 import 체인을 이루며, 모두 동일한 `go_package`를 공유합니다.

```protobuf
option go_package = "example.com/demo/pb";
```

- `base.proto` — 계층 C: 기본 타입(`BaseResp`)을 정의합니다.
- `middleware.proto` — 계층 B: `base.proto`를 import하고 `RequestMeta`를 정의합니다.
- `main.proto` — 계층 A: `middleware.proto`를 import하고 `PingService`(진입점)를 정의합니다.

import 체인: `main.proto` → `middleware.proto` → `base.proto`

## 생성 명령

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 코드를 생성합니다.

```bash
goctl rpc protoc main.proto \
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
│   └── pingsvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── pinglogic.go
│   ├── server
│   │   └── pingserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── base.pb.go
│   ├── main.pb.go
│   ├── main_grpc.pb.go
│   └── middleware.pb.go
├── pingservice
│   └── pingservice.go
└── pingsvc.go
```

## 핵심 사항

- 세 proto 파일(`main.proto` → `middleware.proto` → `base.proto`)이 전이 import 체인을 이룹니다.
- goctl은 모든 전이 import를 자동으로 재귀 해결합니다.
- 세 파일 모두 동일한 `go_package = "example.com/demo/pb"`를 공유합니다.
- 진입 proto 파일만 지정하면 됩니다. 나머지는 goctl과 protoc가 처리합니다.
- 순환 import는 감지되며 오류가 발생합니다(protoc 동작과 동일).
