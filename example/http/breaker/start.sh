#!/bin/bash

GOOS=linux go build -ldflags="-s -w" server.go
docker run --rm -it --cpus=1 -p 8080:8080 -v `pwd`:/app -w /app alpine /app/server
rm -f server
