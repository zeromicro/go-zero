package gen

const (
	dockerTemplate = `FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR $GOPATH/src/{{.projectName}}
COPY . .
RUN go build -ldflags="-s -w" -o /app/{{.exeFile}} {{.goRelPath}}/{{.goFile}}


FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/{{.exeFile}} /app/{{.exeFile}}

CMD ["./{{.exeFile}}"{{.argument}}]
`

	makefileTemplate = `version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	docker pull alpine
	docker pull golang:alpine
	cd $(GOPATH)/src/xiao && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/{{.exeFile}}:$(version) . -f {{.relPath}}/Dockerfile
	docker image prune --filter label=stage=gobuilder -f

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/{{.exeFile}}:$(version)

deploy: push
	kubectl -n {{.namespace}} set image deployment/{{.exeFile}}-deployment {{.exeFile}}=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/{{.exeFile}}:$(version)
`
)
