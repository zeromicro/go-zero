
MAINTAINER=kevinwan
CLI_NAME=goctl
FILE_NAME=goctl
DOCKER_NAME=goctl
VERSION_NAME=$(version)

build:
	go build -ldflags="-s -w" $(FILE_NAME).go
	$(if $(shell command -v upx), upx $(FILE_NAME))

mac:
	GOOS=darwin go build -ldflags="-s -w" -o $(CLI_NAME)-darwin $(FILE_NAME).go
	$(if $(shell command -v upx), upx $(CLI_NAME)-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o $(CLI_NAME).exe $(FILE_NAME).go
	$(if $(shell command -v upx), upx $(CLI_NAME).exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o $(CLI_NAME)-linux $(FILE_NAME).go
	$(if $(shell command -v upx), upx $(CLI_NAME)-linux)

image:
	docker build --rm --platform linux/amd64 -t $(MAINTAINER)/$(DOCKER_NAME):$(version) .
	docker tag $(MAINTAINER)/$(DOCKER_NAME):$(version) $(MAINTAINER)/$(DOCKER_NAME):latest
	docker push $(MAINTAINER)/$(DOCKER_NAME):$(version)
	docker push $(MAINTAINER)/$(DOCKER_NAME):latest
	docker build --rm --platform linux/arm64 -t $(MAINTAINER)/$(DOCKER_NAME):$(version)-arm64 .
	docker tag $(MAINTAINER)/$(DOCKER_NAME):$(version)-arm64 $(MAINTAINER)/$(DOCKER_NAME):latest-arm64
	docker push $(MAINTAINER)/$(DOCKER_NAME):$(version)-arm64
	docker push $(MAINTAINER)/$(DOCKER_NAME):latest-arm64
