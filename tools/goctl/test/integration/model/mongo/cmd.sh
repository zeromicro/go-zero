#!/bin/bash

wd=$(dirname $0)
project=test
testDir=$wd/$project
mkdir -p $testDir

cd $testDir

# go mod init
go mod init $project

# generate cache code
goctl model mongo -t User -c --dir cache
if [ $? -ne 0 ]; then
	exit 1
fi

# generate non-cache code
goctl model mongo -t User --dir nocache
if [ $? -ne 0 ]; then
	exit 1
fi

# go mod tidy
go mod tidy

# code inspection
go test -race ./...
if [ $? -ne 0 ]; then
	echo
fi
