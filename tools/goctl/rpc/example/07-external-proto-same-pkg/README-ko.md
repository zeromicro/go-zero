# 예제 07: 외부 Proto — 동일한 `go_package`

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 두 파일이 **동일한** `go_package`를 공유하는 외부 디렉터리의 proto 파일을 import하는 방법을 보여줍니다.

## Proto 정의

`service.proto`와 `ext.proto`는 모두 동일한 `go_package`를 사용합니다.

```protobuf
option go_package = "example.com/demo/pb";
```

소스 레이아웃:

```
07-external-proto-same-pkg/
├── ext_protos
│   └── ext.proto        # External proto (go_package = "example.com/demo/pb")
├── service.proto        # Service definition (go_package = "example.com/demo/pb")
├── README.md
├── README-cn.md
└── README-ko.md
```

- `ext.proto`는 별도 디렉터리(`ext_protos/`)에 있지만 `service.proto`와 동일한 `go_package`를 가집니다.
- `service.proto`는 `ext.proto`를 import하고 `ext.ExtReq` / `ext.ExtReply`를 RPC 타입으로 사용합니다.

## 생성 명령

먼저 출력 디렉터리에 `go.mod`를 초기화합니다.

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

그런 다음 코드를 생성합니다(`-I ./ext_protos`에 주목하세요).

```bash
goctl rpc protoc service.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . -I ./ext_protos
```

생성되는 디렉터리 구조:

```
output/
├── etc
│   └── svc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── querylogic.go
│   ├── server
│   │   └── queryserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── ext.pb.go
│   ├── service.pb.go
│   └── service_grpc.pb.go
├── queryservice
│   └── queryservice.go
└── svc.go
```

## 핵심 사항

- `ext.proto`는 별도 디렉터리(`ext_protos/`)에 있지만 `service.proto`와 동일한 `go_package`를 가집니다.
- 외부 디렉터리를 proto 검색 경로에 추가하려면 `-I ./ext_protos`를 사용합니다.
- 외부 proto가 **동일한** `go_package`를 가지면 모든 타입이 하나의 Go 패키지로 병합되므로 교차 패키지 import가 필요하지 않습니다.
