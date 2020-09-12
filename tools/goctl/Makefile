version := $(shell /bin/date "+%Y-%m-%d %H:%M")

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" goctl.go
	command -v upx &> /dev/null && upx goctl
mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o goctl-darwin goctl.go
	command -v upx &> /dev/null && upx goctl-darwin
win:
	GOOS=windows go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o goctl.exe goctl.go
	command -v upx &> /dev/null && upx goctl.exe
linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o goctl-linux goctl.go
	command -v upx &> /dev/null && upx goctl-linux
