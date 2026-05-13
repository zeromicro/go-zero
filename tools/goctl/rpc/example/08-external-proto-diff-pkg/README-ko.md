# 예제 08: 외부 proto — 다른 `go_package`

[English](README.md) | [中文](README-cn.md) | 한국어

이 예제는 파일들이 **서로 다른** `go_package` 값을 가져 생성된 Go 코드에 패키지 간 import가 필요한 외부 디렉터리의 proto 파일을 import하는 방법을 보여줍니다.

## proto 정의

proto 파일들은 서로 다른 `go_package` 값을 사용합니다.

- `service.proto`: `go_package = "example.com/demo/pb"`
- `ext_protos/common/types.proto`: `go_package = "example.com/demo/pb/common"`

소스 레이아웃:

```
08-external-proto-diff-pkg/
├── ext_protos
│   └── common
│       └── types.proto    # 외부 proto (go_package = "example.com/demo/pb/common")
├── service.proto          # 서비스 정의 (go_package = "example.com/demo/pb")
├── README.md
├── README-cn.md
└── README-ko.md
```

- `types.proto`는 `go_package = "example.com/demo/pb/common"`을 가지며, 이는 **다른** Go 패키지입니다.
- `service.proto`는 `common.ExtReq` / `common.ExtReply`를 RPC 파라미터 타입으로 직접 사용합니다.

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
├── dataservice
│   └── dataservice.go
├── etc
│   └── svc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── fetchlogic.go
│   ├── server
│   │   └── dataserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── common
│   │   └── types.pb.go
│   ├── service.pb.go
│   └── service_grpc.pb.go
└── svc.go
```

## 핵심 사항

- 외부 proto가 **다른** `go_package`를 가지면 goctl은 패키지 간 Go import를 자동으로 생성합니다.
- goctl은 import된 proto의 `go_package` 옵션을 파싱하여 proto `package` 이름(예: `common`)을 올바른 Go import 경로로 해결합니다.
- `service.proto`는 `common.ExtReq` / `common.ExtReply`를 RPC 파라미터 타입으로 직접 사용합니다.
