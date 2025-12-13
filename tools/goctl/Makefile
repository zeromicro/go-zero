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


#echo $(shell find $(target_dir) -name "*.proto")
#protoc -I example/app/pb --go_out=example/app/pb --go-grpc_out=example/app/pb
.PHONY: generate-proto
generate-proto:
	$(eval target_dir := $(filter-out $@,$(MAKECMDGOALS)))
	@if [ -z "$(target_dir)" ]; then \
		echo "Usage: make generate-proto <target_directory>"; \
		exit 1; \
	fi
	protoc -I $(target_dir) --go_out=$(target_dir) --go-grpc_out=$(target_dir) $(target_dir)/*.proto


