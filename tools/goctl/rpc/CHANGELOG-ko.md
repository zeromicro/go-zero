# 변경 로그

[English](CHANGELOG.md) | [中文](CHANGELOG-cn.md) | 한국어

## 릴리스 예정

### 새로운 기능

#### 외부 proto import 지원(`--proto_path` / `-I`)

`-I` / `--proto_path` 플래그로 외부 디렉터리의 proto 파일을 import할 수 있도록 지원을 추가했으며, 전체 전이 의존성 해결을 제공합니다.

**영향받는 파일:**
- `generator/gen.go` — `ZRpcContext`에 `ProtoPaths` 필드를 추가하고, 코드 생성 전에 `ImportedProtos`를 채우는 `resolveImportedProtos()`를 추가했습니다.
- `generator/genpb.go` — 전이적으로 import된 proto 파일을 자동 탐색해 `protoc` 명령에 추가하는 `buildProtocCmd()`를 추가하고, protoc 출력 경로에 맞는 상대 경로를 계산하는 `relativeToProtoPath()`를 추가했습니다.
- `parser/import.go` — 새 파일(대규모 추가). 재귀적 전이 import 해결을 위한 `ResolveImports()`, import된 proto에서 `go_package` / `package` 메타데이터를 추출하는 `ParseImportedProtos()`, proto `package` 이름으로 O(1) 조회를 제공하는 `BuildProtoPackageMap()`을 구현합니다.
- `parser/proto.go` — `Proto` 구조체에 `ImportedProtos []ImportedProto` 필드를 추가했습니다.
- `cli/cli.go` — `RPCNew`의 `ProtoPaths`를 `ZRpcContext`로 전달합니다.
- `cli/zrpc.go` — `VarStringSliceProtoPath`를 `ZRpcContext.ProtoPaths`로 전달합니다.

**이전 vs 이후:**

| | 이전 | 이후 |
|---|---|---|
| 외부 디렉터리의 proto import | ❌ 지원하지 않음, 모든 타입이 같은 파일에 있어야 함 | ✅ `-I ./ext_protos`로 검색 경로 추가 |
| 전이 import(A → B → C) | ❌ 직접 import만 인식 | ✅ 모든 전이 의존성을 재귀적으로 해결 |
| import된 proto의 `.pb.go` 생성 | ❌ 수동, 파일마다 protoc를 별도로 실행해야 함 | ✅ 자동, import된 proto를 protoc 명령에 추가 |
| proto 검색 경로 | ❌ 원본 파일 디렉터리만 | ✅ protoc와 동일하게 여러 `-I` 경로 지원 |

**동작:**
- proto 파일의 모든 `import` 선언을 전이적으로 순회하며, `google/*` well-known types는 건너뜁니다.
- 각 `-I` 디렉터리에서 import된 파일을 검색하며, 사용자 경로에서 찾을 수 없는 시스템 수준 proto는 조용히 건너뜁니다.
- 발견한 proto 파일을 `protoc` 명령에 추가해 메인 proto와 함께 `.pb.go` 파일이 생성되도록 합니다.

---

#### 패키지 간 타입 해석

import된 proto의 `go_package`가 메인 proto와 **다른** 경우, goctl은 이제 서버, 로직, 클라이언트 코드에서 올바른 Go import 경로와 한정된 타입 참조를 자동 생성합니다.

**영향받는 파일:**
- `generator/typeref.go` — 새 파일. 핵심 타입 해석 엔진:
  - `resolveRPCTypeRef()` — proto RPC 타입(단순 타입, 동일 패키지 점 표기 타입, 패키지 간 점 표기 타입, Google WKT)을 올바른 import 경로를 가진 Go 타입 참조로 해석합니다.
  - `resolveCallTypeRef()` — 타입 alias를 지원하는 클라이언트 코드 생성용 변형입니다.
  - `googleWKTTable` — 지원되는 Google well-known types를 Go 대응 타입으로 매핑하는 테이블입니다.
- `generator/genserver.go` — `genFunctions()`가 이제 요청/응답 타입에 대해 `resolveRPCTypeRef()`를 호출하고 추가 import 경로를 수집합니다.
- `generator/genlogic.go` — `genLogicFunction()`이 `resolveRPCTypeRef()`를 사용합니다. main pb import와 패키지 간 import를 조건부로 포함하는 `addLogicImports()`를 추가했습니다.
- `generator/gencall.go` — `genFunction()`과 `getInterfaceFuncs()`가 타입 alias와 추가 import를 위해 `resolveCallTypeRef()`를 사용합니다. `buildExtraImportLines()` 헬퍼를 추가했습니다.
- `generator/call.tpl` — 패키지 간 import 라인을 위한 `{{.extraImports}}` 플레이스홀더를 추가했습니다.

**이전 vs 이후:**

| Proto 타입 | 이전 | 이후 |
|---|---|---|
| `GetReq`(같은 파일) | `pb.GetReq` | `pb.GetReq`(변경 없음) |
| `ext.ExtReq`(같은 `go_package`) | ❌ 오류: "request type must defined in" | ✅ `pb.ExtReq` — 메인 proto의 Go 패키지에 병합 |
| `common.TypesReq`(다른 `go_package`) | ❌ 오류: "request type must defined in" | ✅ `common.TypesReq` + 자동 생성 `import "example.com/demo/pb/common"` |
| `google.protobuf.Empty` | ❌ 오류: "request type must defined in" | ✅ `emptypb.Empty` + 자동 생성 import |

**동작:**
- 단순 타입(예: `GetReq`)은 추가 import 없이 `pb.GetReq`로 해석됩니다.
- 동일 패키지 점 표기 타입(예: `ext.ExtReq`, `ext`가 같은 `go_package`를 가진 경우)은 `pb.ExtReq`로 해석됩니다.
- 패키지 간 점 표기 타입(예: `common.TypesReq`, `common`이 다른 `go_package`를 가진 경우)은 올바른 Go import 경로가 자동 추가된 `common.TypesReq`로 해석됩니다.

---

#### RPC 파라미터로 Google well-known types 사용

Google protobuf well-known types를 이제 메시지 필드뿐 아니라 RPC 요청/응답 타입으로도 직접 사용할 수 있습니다.

**영향받는 파일:**
- `generator/typeref.go` — `resolveGoogleWKT()` + `googleWKTTable`이 지원되는 표준 타입을 처리합니다.

**이전 vs 이후:**

| Proto 타입 | 이전(RPC 파라미터) | 이후(RPC 파라미터) |
|---|---|---|
| `google.protobuf.Empty` | ❌ 오류 | ✅ `emptypb.Empty` |
| `google.protobuf.Timestamp` | ❌ 오류 | ✅ `timestamppb.Timestamp` |
| `google.protobuf.Duration` | ❌ 오류 | ✅ `durationpb.Duration` |
| `google.protobuf.Any` | ❌ 오류 | ✅ `anypb.Any` |
| `google.protobuf.Struct` | ❌ 오류 | ✅ `structpb.Struct` |
| `google.protobuf.FieldMask` | ❌ 오류 | ✅ `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` | ❌ 오류 | ✅ `wrapperspb.*Value` |

> 참고: 이 타입들은 이전에도 **메시지 필드**로 사용할 수 있었습니다. 이번 변경으로 **RPC 요청/응답 타입**으로도 직접 사용할 수 있습니다.

**지원 타입:**

| Proto 타입 | Go 타입 |
|---|---|
| `google.protobuf.Empty` | `emptypb.Empty` |
| `google.protobuf.Timestamp` | `timestamppb.Timestamp` |
| `google.protobuf.Duration` | `durationpb.Duration` |
| `google.protobuf.Any` | `anypb.Any` |
| `google.protobuf.Struct` | `structpb.Struct` |
| `google.protobuf.Value` | `structpb.Value` |
| `google.protobuf.ListValue` | `structpb.ListValue` |
| `google.protobuf.FieldMask` | `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` (wrappers) | `wrapperspb.*Value` |

---

### 호환성이 깨지는 변경

#### RPC 정의에서 점 표기 타입 이름 허용

이전에는 goctl이 점(.)이 포함된 RPC 요청/응답 타입(예: `base.Req`)을 거부했으며, 모든 타입이 같은 proto 파일에 정의되어 있어야 했습니다. 이 제한이 제거되었습니다.

**이전 vs 이후:**

| Proto 정의 | 이전 | 이후 |
|---|---|---|
| `rpc Fetch(base.Req) returns (base.Reply)` | ❌ 파싱 오류: "request type must defined in xxx.proto" | ✅ 성공적으로 파싱되며, `base.Req`가 import된 proto를 통해 해석됨 |
| `rpc Ping(google.protobuf.Empty) returns (Reply)` | ❌ 파싱 오류: "request type must defined in xxx.proto" | ✅ 성공적으로 파싱되며, `emptypb.Empty`로 해석됨 |

**영향받는 파일:**
- `parser/service.go` — 점(.)이 포함된 타입 이름을 `"request type must defined in"` / `"returns type must defined in"` 오류로 거부하던 검증 루프를 제거했습니다.
- `parser/parser_test.go` — `TestDefaultProtoParseCaseInvalidRequestType`과 `TestDefaultProtoParseCaseInvalidResponseType`을 이름을 변경하고 업데이트하여 점 표기 타입이 이제 성공적으로 파싱되는지 검증합니다.
