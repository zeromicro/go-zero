#!/usr/bin/env sh

.\goctl.go example/app/pb/role.proto --go_out=./example/app/pb --go-grpc_out=./example/app/pb --zrpc_out=./example/app -m --style go_zero -I ./example/app/pb