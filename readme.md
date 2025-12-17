# go-zero
<p align="center">
<img align="center" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">
</p>

go-zero is a web and rpc framework with lots of builtin engineering practices. It’s born to ensure the stability of the busy services with resilience design and has been serving sites with tens of millions of users for years.

<div align=center>

[![codecov](https://codecov.io/gh/zeromicro/go-zero/branch/master/graph/badge.svg)](https://codecov.io/gh/zeromicro/go-zero)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)](https://goreportcard.com/report/github.com/zeromicro/go-zero)
[![Release](https://img.shields.io/github/v/release/zeromicro/go-zero.svg?style=flat-square)](https://github.com/zeromicro/go-zero)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeromicro/go-zero.svg)](https://pkg.go.dev/github.com/zeromicro/go-zero)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Discord](https://img.shields.io/discord/794530774463414292?label=chat&logo=discord)](https://discord.gg/4JQvC5A4Fe)

</div>


## 🤷‍ What is go-zero?
English | [简体中文](readme-cn.md)

<a href="https://trendshift.io/repositories/3263" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3263" alt="zeromicro%2Fgo-zero | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
<a href="https://www.producthunt.com/posts/go-zero?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go&#0045;zero" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=334030&theme=light" alt="go&#0045;zero - A&#0032;web&#0032;&#0038;&#0032;rpc&#0032;framework&#0032;written&#0032;in&#0032;Go&#0046; | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>


go-zero (listed in CNCF Landscape: [https://landscape.cncf.io/?selected=go-zero](https://landscape.cncf.io/?selected=go-zero)) is a web and rpc framework with lots of builtin engineering practices. It’s born to ensure the stability of the busy services with resilience design and has been serving sites with tens of millions of users for years.

go-zero contains simple API description syntax and code generation tool called `goctl`. You can generate Go, iOS, Android, Kotlin, Dart, TypeScript, JavaScript from .api files with `goctl`.

#### Advantages of go-zero:

* Improves the stability of the services with tens of millions of daily active users
* Builtin chained timeout control, concurrency control, rate limit, adaptive circuit breaker, adaptive load shedding, even no configuration needed
* Builtin middlewares also can be integrated into your frameworks
* Simple API syntax, one command to generate a couple of different languages
* Auto validate the request parameters from clients
* Plenty of builtin microservice management and concurrent toolkits

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/architecture-en.png" alt="Architecture" width="1500" />

## Backgrounds of go-zero

In early 2018, we transitioned from a Java+MongoDB monolithic architecture to microservices, choosing:

* **Golang** - High performance, simple syntax, excellent deployment experience, and low resource consumption
* **Self-designed microservice framework** - Better problem isolation, easier feature extension, and faster issue resolution

## Design considerations on go-zero

go-zero follows these core design principles:

* **Simplicity** - Keep it simple, first principle
* **High availability** - Stable under high concurrency
* **Resilience** - Failure-oriented programming with adaptive protection
* **Developer friendly** - Encapsulate complexity, one way to do one thing
* **Easy to extend** - Flexible architecture for growth

## The implementation and features of go-zero

go-zero integrates engineering best practices:

* **Code generation** - Powerful tools to minimize boilerplate
* **Simple API** - Clean interfaces, fully compatible with net/http
* **High performance** - Optimized for speed and efficiency
* **Resilience** - Built-in circuit breaker, rate limiting, load shedding, timeout control
* **Service mesh** - Service discovery, load balancing, call tracing
* **Developer tools** - Auto parameter validation, cache management, metrics and monitoring

![Resilience](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/resilience-en.png)

## Architecture with go-zero

<img width="1067" alt="image" src="https://user-images.githubusercontent.com/1918356/171880372-5010d846-e8b1-4942-8fe2-e2bbb584f762.png">

## Installation

Run the following command under your project:

```shell
go get -u github.com/zeromicro/go-zero
```

## AI-Native Development

The go-zero team provides AI tooling for Claude, GitHub Copilot, Cursor to generate framework-compliant code.

### Three Core Projects

**[ai-context](https://github.com/zeromicro/ai-context)** - Workflow guide for AI assistants

**[zero-skills](https://github.com/zeromicro/zero-skills)** - Pattern library with examples

**[mcp-zero](https://github.com/zeromicro/mcp-zero)** - Code generation tools via Model Context Protocol

### Quick Setup

#### GitHub Copilot
```bash
git submodule add https://github.com/zeromicro/ai-context.git .github/ai-context
ln -s ai-context/00-instructions.md .github/copilot-instructions.md  # macOS/Linux
# Windows: mklink .github\copilot-instructions.md .github\ai-context\00-instructions.md
git submodule update --remote .github/ai-context  # Update
```

#### Cursor
```bash
git submodule add https://github.com/zeromicro/ai-context.git .cursorrules
git submodule update --remote .cursorrules  # Update
```

#### Windsurf
```bash
git submodule add https://github.com/zeromicro/ai-context.git .windsurfrules
git submodule update --remote .windsurfrules  # Update
```

#### Claude Desktop
```bash
git clone https://github.com/zeromicro/mcp-zero.git && cd mcp-zero && go build
# Configure: ~/Library/Application Support/Claude/claude_desktop_config.json
# Or: claude mcp add --transport stdio mcp-zero --env GOCTL_PATH=/path/to/goctl -- /path/to/mcp-zero
```

### How It Works

AI assistants use these tools together:
1. **ai-context** - workflow guidance
2. **zero-skills** - implementation patterns
3. **mcp-zero** - real-time code generation

**Example**: Creating a REST API → AI reads **ai-context** for workflow → calls **mcp-zero** to generate code → references **zero-skills** for patterns → produces production-ready code ✅

## Quick Start

1. Full examples:

     [Rapid development of microservice systems](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)

     [Rapid development of microservice systems - multiple RPCs](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)

2. Install goctl

   ```shell
   # for Go
   go install github.com/zeromicro/go-zero/tools/goctl@latest

   # For Mac
   brew install goctl

   # docker for all platforms
   docker pull kevinwan/goctl
   # run goctl
   docker run --rm -it -v `pwd`:/app kevinwan/goctl --help
   ```

   Ensure goctl is executable and in your $PATH.

3. Create the API file (greet.api):

   ```go
   type (
     Request {
       Name string `path:"name,options=[you,me]"` // parameters are auto validated
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

   Generate .api template:

   ```shell
   goctl api -o greet.api
   ```

4. Generate Go server code

   ```shell
   goctl api go -api greet.api -dir greet
   ```

   Generated structure:

   ```Plain Text
   ├── greet
   │   ├── etc
   │   │   └── greet-api.yaml        // configuration file
   │   ├── greet.go                  // main file
   │   └── internal
   │       ├── config
   │       │   └── config.go         // configuration definition
   │       ├── handler
   │       │   ├── greethandler.go   // get/put/post/delete routes are defined here
   │       │   └── routes.go         // routes list
   │       ├── logic
   │       │   └── greetlogic.go     // request logic can be written here
   │       ├── svc
   │       │   └── servicecontext.go // service context, mysql/redis can be passed in here
   │       └── types
   │           └── types.go          // request/response defined here
   └── greet.api                     // api description file
   ```

   Run the service:

   ```shell
   cd greet
   go mod tidy
   go run greet.go -f etc/greet-api.yaml
   ```

   Default port: 8888 (configurable in etc/greet-api.yaml)

   Test with curl:

   ```shell
   curl -i http://localhost:8888/greet/from/you
   ```

   Response:

   ```http
   HTTP/1.1 200 OK
   Date: Sun, 30 Aug 2020 15:32:35 GMT
   Content-Length: 0
   ```

5. Write business logic

    * Pass dependencies (mysql, redis, etc.) via servicecontext.go
    * Add logic code in the logic package per .api definition

6. Generate client code for multiple languages

   ```shell
   goctl api java -api greet.api -dir greet
   goctl api dart -api greet.api -dir greet
   ...
   ```

## Benchmark

![benchmark](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/benchmark.png)

[Checkout the test code](https://github.com/smallnest/go-web-framework-benchmark)

## Documents

* [Documents](https://go-zero.dev/)
* [Rapid development of microservice systems](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)
* [Rapid development of microservice systems - multiple RPCs](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)
* [Examples](https://github.com/zeromicro/zero-examples)

## Chat group

Join the chat via https://discord.gg/4JQvC5A4Fe

## Cloud Native Landscape

<p float="left">
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-logo.svg" width="200"/>&nbsp;&nbsp;&nbsp;
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-landscape-logo.svg" width="150"/>
</p>

go-zero enlisted in the [CNCF Cloud Native Landscape](https://landscape.cncf.io/?selected=go-zero).

## Give a Star! ⭐

If you like this project or are using it to learn or start your own solution, give it a star to get updates on new releases. Your support matters!
