#!/bin/bash

cd $(dirname $0)

# source functions
source ../../../common/echo.sh

console_tip "mongo  test"

# build goctl
console_step "goctl building"

buildFile=goctl
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $buildFile ../../../../goctl.go
image=goctl-mongo:latest

# docker build
console_step "docker building"
docker build -t $image .
if [ $? -ne 0 ]; then
	rm -f $buildFile
	console_red "docker build failed"
	exit 1
fi

# run docker image
console_step "docker running"
docker run --rm $image
if [ $? -ne 0 ]; then
	rm -f $buildFile
	console_red "docker run failed"
	docker image rm -f $image
	exit 1
fi

rm -f $buildFile
console_green "PASS"
docker image rm -f $image > /dev/null 2>&1
