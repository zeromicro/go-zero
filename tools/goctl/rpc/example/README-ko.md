# RPC 예제

[English / 中文](README.md) | 한국어

이 디렉터리에는 모든 `goctl rpc` 코드 생성 시나리오에 대한 완전한 예제가 포함되어 있습니다.

각 예제에는 다음이 포함됩니다.
- `.proto` 소스 파일
- `README.md`(영어), `README-cn.md`(중국어), `README-ko.md`(한국어) 문서

## 예제

| # | 디렉터리 | 시나리오 | 주요 플래그 |
|---|-----------|----------|-----------|
| 01 | [01-basic](01-basic/) | 기본 단일 서비스, import 없음 | — |
| 02 | [02-import-sibling](02-import-sibling/) | 같은 디렉터리의 proto 파일 import | `--proto_path=.` |
| 03 | [03-import-subdir](03-import-subdir/) | 하위 디렉터리의 proto import | `--proto_path=.` |
| 04 | [04-transitive-import](04-transitive-import/) | 전이 import(A → B → C) | `--proto_path=.` |
| 05 | [05-multiple-services](05-multiple-services/) | 하나의 proto에 여러 서비스 | `--multiple` |
| 06 | [06-wellknown-types](06-wellknown-types/) | 메시지에서 Google well-known types 사용 | — |
| 07 | [07-external-proto-same-pkg](07-external-proto-same-pkg/) | 외부 proto, 동일한 `go_package` | `-I ./ext_protos` |
| 08 | [08-external-proto-diff-pkg](08-external-proto-diff-pkg/) | 외부 proto, 다른 `go_package` | `-I ./ext_protos` |
| 09 | [09-google-types-as-rpc](09-google-types-as-rpc/) | RPC 파라미터로 Google well-known types 사용 | — |
| 10 | [10-streaming](10-streaming/) | 서버/클라이언트/양방향 스트리밍 | — |

## 사전 요구 사항

- [Go](https://go.dev/) 1.22+
- [protoc](https://github.com/protocolbuffers/protobuf/releases)(Protocol Buffers 컴파일러)
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) 및 [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)
- [goctl](https://github.com/zeromicro/go-zero/tree/master/tools/goctl)

## 빠른 시작

```bash
# protoc 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 기본 예제 실행
cd 01-basic
mkdir -p output && cd output && go mod init example.com/demo && cd ..
goctl rpc protoc greeter.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```
