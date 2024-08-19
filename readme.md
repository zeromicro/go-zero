# go-zero
<p align="center">
<img align="center" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">
</p>

go-zero is a web and rpc framework with lots of builtin engineering practices. It’s born to ensure the stability of the busy services with resilience design and has been serving sites with tens of millions of users for years.

<div align=center>

[![Go](https://github.com/zeromicro/go-zero/workflows/Go/badge.svg?branch=master)](https://github.com/zeromicro/go-zero/actions)
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

## Design considerations on go-zero

By designing the microservice architecture, we expected to ensure stability, as well as productivity. And from just the beginning, we have the following design principles:

* Keep it simple
* High availability
* Stable on high concurrency
* Easy to extend
* Resilience design, failure-oriented programming
* Try best to be friendly to the business logic development, encapsulate the complexity
* One thing, one way

After almost half a year, we finished the transfer from a monolithic system to microservice system and deployed on August 2018. The new system guaranteed business growth and system stability.

## The implementation and features of go-zero

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

##  The simplified architecture that we use with go-zero

<img width="1067" alt="image" src="https://user-images.githubusercontent.com/1918356/171880372-5010d846-e8b1-4942-8fe2-e2bbb584f762.png">

##  Installation

Run the following command under your project:

```shell
go get -u github.com/zeromicro/go-zero
```

##  Quick Start

1. Full examples can be checked out from below:

     [Rapid development of microservice systems](https://github.com/zeromicro/zero-doc/blob/main/doc/shorturl-en.md)

     [Rapid development of microservice systems - multiple RPCs](https://github.com/zeromicro/zero-doc/blob/main/docs/zero/bookstore-en.md)

2. Install goctl

   `goctl`can be read as `go control`. `goctl` means not to be controlled by code, instead, we control it. The inside `go` is not `golang`. At the very beginning, I was expecting it to help us improve productivity, and make our lives easier.

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
   
   make sure goctl is executable and in your $PATH.
   
3. Create the API file, like greet.api, you can install the plugin of goctl in vs code, api syntax is supported.

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

   the .api files also can be generated by goctl, like below:

   ```shell
   goctl api -o greet.api
   ```

4. Generate the go server-side code

   ```shell
   goctl api go -api greet.api -dir greet
   ```

   the generated files look like:

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

   the generated code can be run directly:

   ```shell
   cd greet
   go mod tidy
   go run greet.go -f etc/greet-api.yaml
   ```

   by default, it’s listening on port 8888, while it can be changed in the configuration file.

   you can check it by curl:

   ```shell
   curl -i http://localhost:8888/greet/from/you
   ```

   the response looks like below:

   ```http
   HTTP/1.1 200 OK
   Date: Sun, 30 Aug 2020 15:32:35 GMT
   Content-Length: 0
   ```

5. Write the business logic code

    * the dependencies can be passed into the logic within servicecontext.go, like mysql, redis, etc.
    * add the logic code in a logic package according to .api file

6. Generate code like Java, TypeScript, Dart, JavaScript, etc. just from the api file

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

##  Chat group

Join the chat via https://discord.gg/4JQvC5A4Fe

##  Cloud Native Landscape

<p float="left">
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-logo.svg" width="200"/>&nbsp;&nbsp;&nbsp;
<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/cncf-landscape-logo.svg" width="150"/>
</p>

go-zero enlisted in the [CNCF Cloud Native Landscape](https://landscape.cncf.io/?selected=go-zero).

## Give a Star! ⭐

If you like this project or are using it to learn or start your own solution, give it a star to get updates on new releases. Your support matters!

## Buy me a coffee

<a href="https://www.buymeacoffee.com/kevwan" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
