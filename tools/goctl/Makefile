build:
	go build -ldflags="-s -w" goctl.go
	$(if $(shell command -v upx || which upx), upx goctl)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o goctl-darwin goctl.go
	$(if $(shell command -v upx || which upx), upx goctl-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o goctl.exe goctl.go
	$(if $(shell command -v upx || which upx), upx goctl.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o goctl-linux goctl.go
	$(if $(shell command -v upx || which upx), upx goctl-linux)

image:
	docker build --rm --platform linux/amd64 -t kevinwan/goctl:$(version) .
	docker tag kevinwan/goctl:$(version) kevinwan/goctl:latest
	docker push kevinwan/goctl:$(version)
	docker push kevinwan/goctl:latest
	docker build --rm --platform linux/arm64 -t kevinwan/goctl:$(version)-arm64 .
	docker tag kevinwan/goctl:$(version)-arm64 kevinwan/goctl:latest-arm64
	docker push kevinwan/goctl:$(version)-arm64
	docker push kevinwan/goctl:latest-arm64
