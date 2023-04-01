GOPATH?=$(shell go env GOPATH)
build:
	go build -ldflags="-s -w" goctl.go
	mv goctl $(GOPATH)/bin/goctls
	goctls template clean
	goctls template init
	$(if $(shell command -v upx), upx goctl)

mac:
	@echo ${GOPATH}
	GOOS=darwin go build -ldflags="-s -w" -o goctls goctl.go
	$(if $(shell command -v upx), upx goctl-darwin)
	mv goctls $(GOPATH)/bin/goctls

win:
	GOOS=windows go build -ldflags="-s -w" -o goctls.exe goctl.go
	$(if $(shell command -v upx), upx goctl.exe)
	mv goctls.exe $(GOPATH)/bin/goctls.exe

linux:
	GOOS=linux go build -ldflags="-s -w" -o goctls goctl.go
	$(if $(shell command -v upx), upx goctl-linux)
	mv goctls $(GOPATH)/bin/goctls

image:
	docker build --rm --platform linux/amd64 -t kevinwan/goctl:$(version) .
	docker tag kevinwan/goctl:$(version) kevinwan/goctl:latest
	docker push kevinwan/goctl:$(version)
	docker push kevinwan/goctl:latest
	docker build --rm --platform linux/arm64 -t kevinwan/goctl:$(version)-arm64 .
	docker tag kevinwan/goctl:$(version)-arm64 kevinwan/goctl:latest-arm64
	docker push kevinwan/goctl:$(version)-arm64
	docker push kevinwan/goctl:latest-arm64
