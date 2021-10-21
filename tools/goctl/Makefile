build:
	go build -ldflags="-s -w" goctl.go
	$(if $(shell command -v upx), upx goctl)
mac:
	GOOS=darwin go build -ldflags="-s -w" -o goctl-darwin goctl.go
	$(if $(shell command -v upx), upx goctl-darwin)
win:
	GOOS=windows go build -ldflags="-s -w" -o goctl.exe goctl.go
	$(if $(shell command -v upx), upx goctl.exe)
linux:
	GOOS=linux go build -ldflags="-s -w" -o goctl-linux goctl.go
	$(if $(shell command -v upx), upx goctl-linux)
