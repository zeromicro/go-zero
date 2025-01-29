#!/bin/bash

wd=$(pwd)
output="$wd/hi"

rm -rf $output

goctl rpc protoc -I $wd "$wd/hi.proto" --go_out="$output/pb" --go-grpc_out="$output/pb" --zrpc_out="$output" --multiple

if [ $? -ne 0 ]; then
    echo "Generate failed"
    exit 1
fi

go mod tidy

if [ $? -ne 0 ]; then
    echo "Tidy failed"
    exit 1
fi

go test ./...