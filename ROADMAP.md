# go-zero Roadmap

This document defines a high level roadmap for go-zero development and upcoming releases.
Community and contributor involvement is vital for successfully implementing all desired items for each release.
We hope that the items listed below will inspire further engagement from the community to keep go-zero progressing and shipping exciting and valuable features.

## 2021 Q2
- Support TLS in redis connections
- Support service discovery through K8S watch api
- Log full sql statements for easier sql problem solving

## 2021 Q3
- Support `goctl mock` command to start a mocking server with given `.api` file
- Adapt builtin tracing mechanism to opentracing solutions
- Support `goctl model pg` to support PostgreSQL code generation

## 2021 Q4
- Support `goctl doctor` command to report potential issues for given service
- Support `context` in redis related methods for timeout and tracing
- Support `context` in sql related methods for timeout and tracing
- Support `context` in mongodb related methods for timeout and tracing