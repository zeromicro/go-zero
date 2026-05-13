# 예제 03: 하위 디렉터리의 proto import

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 하위 디렉터리에 있는 proto 파일을 import하는 방법을 보여줍니다.

## proto 정의

두 proto 파일은 **서로 다른** `go_package` 값을 사용합니다.

- `order.proto` — `OrderService`를 정의하고 `common/types.proto`를 import합니다.

```protobuf
option go_package = "example.com/demo/pb";
```

- `common/types.proto` — 재사용 가능한 페이지네이션 및 정렬 메시지를 정의합니다.

```protobuf
option go_package = "example.com/demo/pb/common";
```

`order.proto`는 하위 디렉터리에서 `common/types.proto`를 import합니다.

```protobuf
import "common/types.proto";
```

두 파일은 **서로 다른** `go_package` 값을 가지므로 별도의 Go 패키지로 컴파일됩니다.

## 생성 명령

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 코드를 생성합니다.

```bash
goctl rpc protoc order.proto \
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
│   └── ordersvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── getorderlogic.go
│   │   └── listorderslogic.go
│   ├── server
│   │   └── orderserviceserver.go
│   └── svc
│       └── servicecontext.go
├── orderservice
│   └── orderservice.go
├── ordersvc.go
└── pb
    ├── common
    │   └── types.pb.go
    ├── order.pb.go
    └── order_grpc.pb.go
```

## 핵심 사항

- 두 proto 파일은 **서로 다른** `go_package` 값을 가지므로 별도의 Go 패키지(`pb/`와 `pb/common/`)로 컴파일됩니다.
- `order.proto`는 하위 디렉터리에서 `common/types.proto`를 import합니다.
- import된 proto의 `go_package`가 다르면 goctl은 교차 패키지 import를 자동으로 생성합니다.
- `-I .` 플래그는 protoc에 현재 디렉터리부터 검색하라고 알려주어 `common/types.proto`를 찾을 수 있게 합니다.
