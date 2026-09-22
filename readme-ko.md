# go-zero
<p align="center">
<img align="center" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">
</p>

go-zero는 다양한 엔지니어링 모범 사례가 내장된 웹 및 RPC 프레임워크입니다. 회복 탄력성 설계로 트래픽이 많은 서비스의 안정성을 보장하기 위해 만들어졌으며, 수년 동안 수천만 사용자를 보유한 사이트에서 운영되어 왔습니다.

<div align=center>

[![codecov](https://codecov.io/gh/zeromicro/go-zero/branch/master/graph/badge.svg)](https://codecov.io/gh/zeromicro/go-zero)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)](https://goreportcard.com/report/github.com/zeromicro/go-zero)
[![Release](https://img.shields.io/github/v/release/zeromicro/go-zero.svg?style=flat-square)](https://github.com/zeromicro/go-zero)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeromicro/go-zero.svg)](https://pkg.go.dev/github.com/zeromicro/go-zero)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Discord](https://img.shields.io/discord/794530774463414292?label=chat&logo=discord)](https://discord.gg/4JQvC5A4Fe)

</div>


## 🤷‍ go-zero란?
[English](readme.md) | [简体中文](readme-cn.md) | 한국어

<a href="https://trendshift.io/repositories/3263" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3263" alt="zeromicro%2Fgo-zero | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
<a href="https://www.producthunt.com/posts/go-zero?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go&#0045;zero" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=334030&theme=light" alt="go&#0045;zero - A&#0032;web&#0032;&#0038;&#0032;rpc&#0032;framework&#0032;written&#0032;in&#0032;Go&#0046; | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>


go-zero는 [CNCF Landscape](https://landscape.cncf.io/?selected=go-zero)에 등재된, 다양한 엔지니어링 모범 사례가 내장된 웹 및 RPC 프레임워크입니다. 회복 탄력성 설계로 트래픽이 많은 서비스의 안정성을 보장하기 위해 만들어졌으며, 수년 동안 수천만 사용자를 보유한 사이트에서 운영되어 왔습니다.

go-zero에는 간단한 API 설명 문법과 `goctl`이라는 코드 생성 도구가 포함되어 있습니다. `goctl`을 사용하면 .api 파일에서 Go, iOS, Android, Kotlin, Dart, TypeScript, JavaScript 코드를 생성할 수 있습니다.

#### go-zero의 장점:

* 수천만 일간 활성 사용자를 처리하는 서비스의 안정성을 높입니다.
* 체인형 타임아웃 제어, 동시성 제어, 속도 제한, 적응형 서킷 브레이커, 적응형 로드 셰딩이 내장되어 있으며, 별도 설정 없이도 사용할 수 있습니다.
* 내장 미들웨어를 기존 프레임워크에 통합할 수 있습니다.
* 간단한 API 문법과 한 번의 명령으로 여러 언어 코드를 생성할 수 있습니다.
* 클라이언트 요청 파라미터를 자동으로 검증합니다.
* 다양한 내장 마이크로서비스 관리 도구와 동시성 도구를 제공합니다.

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/architecture-en.png" alt="Architecture" width="1500" />

## go-zero의 배경

2018년 초, 우리는 Java+MongoDB 기반 모놀리식 아키텍처에서 마이크로서비스로 전환하면서 다음을 선택했습니다.

* **Golang** - 높은 성능, 간결한 문법, 뛰어난 배포 경험, 낮은 리소스 사용량
* **자체 설계한 마이크로서비스 프레임워크** - 더 나은 문제 격리, 더 쉬운 기능 확장, 더 빠른 이슈 해결

## go-zero의 설계 고려 사항

go-zero는 다음 핵심 설계 원칙을 따릅니다.

* **단순성** - 단순하게 유지하는 것을 제1원칙으로 삼습니다.
* **고가용성** - 높은 동시성 상황에서도 안정적으로 동작합니다.
* **회복 탄력성** - 적응형 보호를 갖춘 장애 지향 프로그래밍을 지향합니다.
* **개발자 친화성** - 복잡성을 캡슐화하고 한 가지 일을 하는 한 가지 방법을 제공합니다.
* **확장 용이성** - 성장을 위한 유연한 아키텍처를 제공합니다.

## go-zero의 구현과 기능

go-zero는 엔지니어링 모범 사례를 통합합니다.

* **코드 생성** - 보일러플레이트를 최소화하는 강력한 도구
* **간단한 API** - 깔끔한 인터페이스와 net/http 완전 호환
* **고성능** - 속도와 효율성을 위한 최적화
* **회복 탄력성** - 내장 서킷 브레이커, 속도 제한, 로드 셰딩, 타임아웃 제어
* **서비스 메시** - 서비스 디스커버리, 부하 분산, 호출 추적
* **개발자 도구** - 자동 파라미터 검증, 캐시 관리, 메트릭과 모니터링

![Resilience](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/resilience-en.png)

## go-zero 아키텍처

<img width="1067" alt="image" src="https://user-images.githubusercontent.com/1918356/171880372-5010d846-e8b1-4942-8fe2-e2bbb584f762.png">

## 설치

프로젝트에서 다음 명령을 실행하세요.

```shell
go get -u github.com/zeromicro/go-zero
```

## AI 네이티브 개발

go-zero 팀은 Claude Code, GitHub Copilot, Cursor 등의 AI 코딩 어시스턴트를 위한 워크플로 가이드와 구현 패턴을 제공합니다. 터미널을 사용할 수 있는 어시스턴트에는 `goctl`을 직접 실행하여 프레임워크 규칙을 따르는 코드를 생성하는 방식을 권장합니다.

### AI 도구 프로젝트

- **[ai-context](https://github.com/zeromicro/ai-context)** - go-zero 개발 작업의 진행 방식과 goctl 사용 시점을 안내하는 워크플로 가이드.
- **[zero-skills](https://github.com/zeromicro/zero-skills)** - 필요할 때 참조하는 구현 패턴, 예제, goctl 명령어 레퍼런스.
- **[mcp-zero](https://github.com/zeromicro/mcp-zero)** - Model Context Protocol(MCP)을 통해 goctl 코드 생성 기능을 제공하는 선택적 어댑터.

### 빠른 설정

goctl을 설치하고 `$PATH`에 포함되어 있는지 확인하세요.

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
goctl --version
```

[ai-context 설정 가이드](https://github.com/zeromicro/ai-context#manual-setup)와 [zero-skills 통합 가이드](https://github.com/zeromicro/zero-skills/tree/main/getting-started)를 참고하여 어시스턴트의 워크플로 지침과 패턴 라이브러리를 설정하세요.

### 동작 방식

1. **ai-context**에서 워크플로를 확인하고 **zero-skills**에서 관련 패턴을 참조합니다.
2. `.api` 또는 `.proto` 명세를 작성한 뒤 터미널에서 **goctl**을 실행하여 코드를 생성합니다.
3. 비즈니스 로직을 구현하고 의존성을 정리한 뒤 서비스를 빌드하고 테스트합니다.

**예시**: REST API 생성 → `.api` 명세 작성 → `goctl api go` 실행 → **zero-skills** 패턴을 참고하여 로직 구현 → `go mod tidy`, `go build ./...`, `go test ./...` 실행.

### 선택적 MCP 통합

Claude Desktop처럼 MCP 도구를 사용하는 클라이언트에서는 **[mcp-zero](https://github.com/zeromicro/mcp-zero)**를 설정하여 MCP를 통해 goctl을 호출할 수 있습니다. 두 방식은 같은 코드 생성 엔진을 사용하며, 터미널을 사용할 수 있는 어시스턴트는 MCP 서버 없이 위의 직접 실행 워크플로를 사용할 수 있습니다.

## 빠른 시작

1. 전체 예제:

     [마이크로서비스 시스템 빠른 개발](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)

     [마이크로서비스 시스템 빠른 개발 - 다중 RPC](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)

2. goctl 설치

   ```shell
   # Go 사용
   go install github.com/zeromicro/go-zero/tools/goctl@latest

   # Mac 사용
   brew install goctl

   # 모든 플랫폼에서 docker 사용
   docker pull kevinwan/goctl
   # goctl 실행
   docker run --rm -it -v `pwd`:/app kevinwan/goctl --help
   ```

   goctl이 실행 가능하고 $PATH에 포함되어 있는지 확인하세요.

3. API 파일(greet.api) 생성:

   ```go
   type (
     Request {
       Name string `path:"name,options=[you,me]"` // 파라미터는 자동으로 검증됩니다.
     }

     Response {
       Message string `json:"message"`
     }
   )

   service greet-api {
     @handler GreetHandler
     get /greet/from/:name(Request) returns (Response)
   }
   ```

   .api 템플릿 생성:

   ```shell
   goctl api -o greet.api
   ```

4. Go 서버 코드 생성

   ```shell
   goctl api go -api greet.api -dir greet
   ```

   생성되는 구조:

   ```Plain Text
   ├── greet
   │   ├── etc
   │   │   └── greet-api.yaml        // 설정 파일
   │   ├── greet.go                  // main 파일
   │   └── internal
   │       ├── config
   │       │   └── config.go         // 설정 정의
   │       ├── handler
   │       │   ├── greethandler.go   // get/put/post/delete 라우트가 여기에 정의됩니다.
   │       │   └── routes.go         // 라우트 목록
   │       ├── logic
   │       │   └── greetlogic.go     // 요청 로직을 여기에 작성할 수 있습니다.
   │       ├── svc
   │       │   └── servicecontext.go // 서비스 컨텍스트, mysql/redis를 여기에 전달할 수 있습니다.
   │       └── types
   │           └── types.go          // 요청/응답이 여기에 정의됩니다.
   └── greet.api                     // api 설명 파일
   ```

   서비스 실행:

   ```shell
   cd greet
   go mod tidy
   go run greet.go -f etc/greet-api.yaml
   ```

   기본 포트: 8888(etc/greet-api.yaml에서 변경 가능)

   curl로 테스트:

   ```shell
   curl -i http://localhost:8888/greet/from/you
   ```

   응답:

   ```http
   HTTP/1.1 200 OK
   Date: Sun, 30 Aug 2020 15:32:35 GMT
   Content-Length: 0
   ```

5. 비즈니스 로직 작성

    * servicecontext.go를 통해 의존성(mysql, redis 등)을 전달합니다.
    * .api 정의에 따라 logic 패키지에 로직 코드를 추가합니다.

6. 여러 언어의 클라이언트 코드 생성

   ```shell
   goctl api java -api greet.api -dir greet
   goctl api dart -api greet.api -dir greet
   ...
   ```

## 벤치마크

![benchmark](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/benchmark.png)

[테스트 코드 보기](https://github.com/smallnest/go-web-framework-benchmark)

## 문서

* [문서](https://go-zero.dev/)
* [마이크로서비스 시스템 빠른 개발](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)
* [마이크로서비스 시스템 빠른 개발 - 다중 RPC](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)
* [예제](https://github.com/zeromicro/zero-examples)

## 채팅 그룹

다음 링크에서 채팅에 참여하세요: https://discord.gg/4JQvC5A4Fe

## Cloud Native Landscape

<p float="left">
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-logo.svg" width="200"/>&nbsp;&nbsp;&nbsp;
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-landscape-logo.svg" width="150"/>
</p>

go-zero는 [CNCF Cloud Native Landscape](https://landscape.cncf.io/?selected=go-zero)에 등재되었습니다.

## 별을 눌러주세요! ⭐

이 프로젝트가 마음에 들거나 학습 또는 자체 솔루션을 시작하는 데 사용 중이라면, 새 릴리스 업데이트를 받을 수 있도록 star를 눌러주세요. 여러분의 지원은 큰 힘이 됩니다!
