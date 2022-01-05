#! /bin/sh

starttest() {
    set -e
    GO111MODULE=on go test -race ./...
}

if [ -z "${TEST_MOD}" ]; then
    docker run --rm --name go-zero -ti \
    -v `pwd`:/go/src/github.com/zeromicro/go-zero \
    --workdir /go/src/github.com/zeromicro/go-zero \
    --env GOPROXY=https://goproxy.cn\
    --env TEST_MOD=local \
    golang:1.16 \
    sh test.sh
else
    starttest
fi

