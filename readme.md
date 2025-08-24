# go-zero
<p align="center">
<img align="center" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">
</p>

go-zero is a web and rpc framework with lots of builtin engineering practices. It's born to ensure the stability of the busy services with resilience design and has been serving sites with tens of millions of users for years.

<div align="center">

[![codecov](https://codecov.io/gh/zeromicro/go-zero/branch/master/graph/badge.svg)](https://codecov.io/gh/zeromicro/go-zero)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)](https://goreportcard.com/report/github.com/zeromicro/go-zero)
[![Release](https://img.shields.io/github/v/release/zeromicro/go-zero.svg?style=flat-square)](https://github.com/zeromicro/go-zero)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeromicro/go-zero.svg)](https://pkg.go.dev/github.com/zeromicro/go-zero)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Discord](https://img.shields.io/discord/794530774463414292?label=chat&logo=discord)](https://discord.gg/4JQvC5A4Fe)

</div>

## 📋 Table of Contents
- [🤷‍ What is go-zero?](#-what-is-go-zero)
- [🚀 Key Features](#-key-features)
- [📐 Architecture](#-architecture)
- [🎯 Background](#-background)
- [🛠 Design Considerations](#-design-considerations)
- [⚡ Implementation and Features](#-implementation-and-features)
- [🏗 Our Architecture](#-our-architecture)
- [📦 Installation](#-installation)
- [🚀 Quick Start](#-quick-start)
- [📊 Benchmark](#-benchmark)
- [📚 Documentation](#-documentation)
- [💬 Community](#-community)
- [☁️ Cloud Native](#-cloud-native)
- [⭐ Give a Star!](#-give-a-star)

## 🤷‍ What is go-zero?
English | [简体中文](readme-cn.md)

<a href="https://trendshift.io/repositories/3263" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3263" alt="zeromicro%2Fgo-zero | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
<a href="https://www.producthunt.com/posts/go-zero?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go&#0045;zero" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=334030&theme=light" alt="go&#0045;zero - A&#0032;web&#0032;&#0038;&#0032;rpc&#0032;framework&#0032;written&#0032;in&#0032;Go&#0046; | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

go-zero (listed in [CNCF Landscape](https://landscape.cncf.io/?selected=go-zero)) is a web and rpc framework with lots of builtin engineering practices. It's born to ensure the stability of the busy services with resilience design and has been serving sites with tens of millions of users for years.

go-zero contains simple API description syntax and code generation tool called `goctl`. You can generate Go, iOS, Android, Kotlin, Dart, TypeScript, JavaScript from .api files with `goctl`.

## 🚀 Key Features

### 🎯 **Production-Ready Stability**
- ✅ **Battle-tested** - Serving tens of millions of daily active users
- ✅ **Resilience design** - Built-in failure-oriented programming
- ✅ **High availability** - Proven in high-traffic production environments

### ⚡ **Built-in Microservice Governance**
- 🔄 **Automatic circuit breaker** - Self-triggering and self-recovering
- 🚦 **Rate limiting** - Adaptive concurrency control
- 📉 **Load shedding** - Intelligent traffic management
- ⏱️ **Timeout control** - Cascading timeout management
- 🔍 **Service discovery** - Built-in service registry and discovery

### 🛠️ **Developer Experience**
- 📝 **Simple API syntax** - Minimal, intuitive API definitions
- 🎨 **Code generation** - One command generates multiple languages
- ✅ **Auto validation** - Request parameter validation out-of-the-box
- 🔧 **Powerful tooling** - Comprehensive `goctl` CLI tool
- 📦 **Middleware support** - Easy to extend and customize

### 🌐 **Multi-Language Support**
Generate client code for: **Go** • **Java** • **TypeScript** • **JavaScript** • **Dart** • **Kotlin** • **iOS** • **Android**

## 📐 Architecture

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/architecture-en.png" alt="Architecture" width="1500" />

## 🎯 Background

In early 2018, we embarked on a transformative journey to redesign our system, transitioning from a monolithic architecture built with Java and MongoDB to a microservices architecture. After careful research and comparison, we made a deliberate choice to:

* Go Beyond with Golang
  * Great performance
  * Simple syntax
  * Proven engineering efficiency
  * Extreme deployment experience
  * Less server resource consumption

* Self-Design Our Microservice Architecture
  * Microservice architecture facilitates the creation of scalable, flexible, and maintainable software systems with independent, reusable components.
  * Easy to locate the problems within microservices.
  * Easy to extend the features by adding or modifying specific microservices without impacting the entire system.

## 🛠 Design Considerations

By designing the microservice architecture, we expected to ensure stability, as well as productivity. And from just the beginning, we have the following design principles:

* Keep it simple
* High availability
* Stable on high concurrency
* Easy to extend
* Resilience design, failure-oriented programming
* Try best to be friendly to the business logic development, encapsulate the complexity
* One thing, one way

After almost half a year, we finished the transfer from a monolithic system to microservice system and deployed on August 2018. The new system guaranteed business growth and system stability.

## ⚡ Implementation and Features

go-zero is a web and rpc framework that integrates lots of engineering practices. The features are mainly listed below:

* Powerful tool included, less code to write
* Simple interfaces
* Fully compatible with net/http
* Middlewares are supported, easy to extend
* High performance
* Failure-oriented programming, resilience design
* Builtin service discovery, load balancing
* Builtin concurrency control, adaptive circuit breaker, adaptive load shedding, auto-trigger, auto recover
* Auto validation of API request parameters
* Chained timeout control
* Auto management of data caching
* Call tracing, metrics, and monitoring
* High concurrency protected

As below, go-zero protects the system with a couple of layers and mechanisms:

![Resilience](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/resilience-en.png)

## 🏗 Our Architecture

<img width="1067" alt="image" src="https://user-images.githubusercontent.com/1918356/171880372-5010d846-e8b1-4942-8fe2-e2bbb584f762.png">

## 📦 Installation

Run the following command under your project:

```shell
go get -u github.com/zeromicro/go-zero
```

## 🚀 Quick Start

### 🔧 Prerequisites
- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Basic Go knowledge** - Understanding of Go syntax and modules

### 📚 Full Examples
- 📖 [Rapid development of microservice systems](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)
- 🔗 [Multiple RPCs example](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)

### ⚡ 5-Minute Setup

#### Step 1: Install goctl

`goctl` is the Swiss Army knife for go-zero development.

```bash
# For Go
go install github.com/zeromicro/go-zero/tools/goctl@latest

# For Mac users
brew install goctl

# Using Docker (all platforms)
docker pull kevinwan/goctl
# Run goctl
docker run --rm -it -v `pwd`:/app kevinwan/goctl --help
```

#### Step 2: Create Your First API

Create a simple API definition file `greet.api`:

```go
type (
  Request {
    Name string `path:"name,options=[you,me]"` // auto-validated parameters
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

#### Step 3: Generate Code

Generate a complete service with one command:

```bash
goctl api go -api greet.api -dir greet
```

This creates a production-ready structure:

```
├── greet
│   ├── etc
│   │   └── greet-api.yaml        // ⚙️  Configuration
│   ├── greet.go                  // 🚀 Main entry point
│   └── internal
│       ├── config
│       │   └── config.go         // 📝 Config definitions
│       ├── handler
│       │   ├── greethandler.go   // 🎯 HTTP handlers
│       │   └── routes.go         // 🛣️  Route definitions
│       ├── logic
│       │   └── greetlogic.go     // 💼 Business logic
│       ├── svc
│       │   └── servicecontext.go // 🔧 Service context
│       └── types
│           └── types.go          // 📊 Type definitions
```

#### Step 4: Run Your Service

```bash
cd greet
go mod tidy
go run greet.go -f etc/greet-api.yaml
```

Your service is now running on `http://localhost:8888`! 🎉

#### Step 5: Test It

```bash
curl -i http://localhost:8888/greet/from/you
```

Response:
```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 01 Jan 2024 12:00:00 GMT
Content-Length: 25

{"message":"Hello you!"}
```

### 🎯 Next Steps

1. **Add Business Logic** - Edit `internal/logic/greetlogic.go`
2. **Add Dependencies** - Configure database, Redis, etc. in `servicecontext.go`
3. **Generate Clients** - Use `goctl` to generate client code for other languages

```bash
# Generate TypeScript client
goctl api ts -api greet.api -dir ./ts-client

# Generate Java client
goctl api java -api greet.api -dir ./java-client

# Generate Dart client
goctl api dart -api greet.api -dir ./dart-client
```

### 📚 More Examples

#### 🏗️ **API Service Examples**
- **[Simple CRUD API](https://github.com/zeromicro/zero-examples/tree/main/http)** - Basic REST API with CRUD operations
- **[JWT Authentication](https://github.com/zeromicro/zero-examples/tree/main/http/jwt)** - API with JWT token authentication
- **[File Upload API](https://github.com/zeromicro/zero-examples/tree/main/http/fileupload)** - Handle file uploads and downloads

#### ⚡ **RPC Service Examples**
- **[Simple RPC](https://github.com/zeromicro/zero-examples/tree/main/rpc/simple)** - Basic gRPC service
- **[RPC with Database](https://github.com/zeromicro/zero-examples/tree/main/rpc/crud)** - RPC service with MySQL integration
- **[Stream RPC](https://github.com/zeromicro/zero-examples/tree/main/rpc/stream)** - Server and client streaming

#### 🌐 **Microservice Examples**
- **[API Gateway](https://github.com/zeromicro/zero-examples/tree/main/gateway)** - API Gateway with multiple backend services
- **[Service Discovery](https://github.com/zeromicro/zero-examples/tree/main/discovery)** - Service registration and discovery
- **[Circuit Breaker Demo](https://github.com/zeromicro/zero-examples/tree/main/breaker)** - Resilience patterns in action

#### 🗄️ **Database Integration**
- **[MySQL Example](https://github.com/zeromicro/zero-examples/tree/main/mysql)** - MySQL database operations
- **[Redis Cache](https://github.com/zeromicro/zero-examples/tree/main/redis)** - Redis caching implementation
- **[MongoDB Example](https://github.com/zeromicro/zero-examples/tree/main/mongo)** - MongoDB document operations

## 📊 Benchmark

![benchmark](https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/benchmark.png)

[Checkout the test code](https://github.com/smallnest/go-web-framework-benchmark)

## 📚 Documentation

* [📖 Official Documentation](https://go-zero.dev/)
* [🚀 Rapid development of microservice systems](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)
* [🔗 Multiple RPCs example](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)
* [💡 Examples Repository](https://github.com/zeromicro/zero-examples)

## 💬 Community

Join the chat via https://discord.gg/4JQvC5A4Fe

## ☁️ Cloud Native

<p float="left">
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-logo.svg" width="200"/>&nbsp;&nbsp;&nbsp;
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-landscape-logo.svg" width="150"/>
</p>

go-zero enlisted in the [CNCF Cloud Native Landscape](https://landscape.cncf.io/?selected=go-zero).

## ⭐ Give a Star!

If you like this project or are using it to learn or start your own solution, give it a star to get updates on new releases. Your support matters!
