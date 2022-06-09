# go-zero Roadmap

This document defines a high level roadmap for go-zero development and upcoming releases.
Community and contributor involvement is vital for successfully implementing all desired items for each release.
We hope that the items listed below will inspire further engagement from the community to keep go-zero progressing and shipping exciting and valuable features.

## 2021 Q2
- [x] Support service discovery through K8S client api
- [x] Log full sql statements for easier sql problem solving

## 2021 Q3
- [x] Support `goctl model pg` to support PostgreSQL code generation
- [x] Adapt builtin tracing mechanism to opentracing solutions

## 2021 Q4
- [x] Support `username/password` authentication in ETCD
- [x] Support `SSL/TLS` in ETCD
- [x] Support `SSL/TLS` in `zRPC`
- [x] Support `TLS` in redis connections
- [x] Support `goctl bug` to report bugs conveniently

## 2022
- [x] Support `context` in redis related methods for timeout and tracing
- [x] Support `context` in sql related methods for timeout and tracing
- [x] Support `context` in mongodb related methods for timeout and tracing
- [x] Add `httpc.Do` with HTTP call governance, like circuit breaker etc.
- [ ] Support `goctl doctor` command to report potential issues for given service
- [ ] Support `goctl mock` command to start a mocking server with given `.api` file
