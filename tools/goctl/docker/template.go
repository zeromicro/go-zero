package docker

import (
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const (
	category           = "docker"
	dockerTemplateFile = "docker.tpl"
	dockerTemplate     = `FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/zero
COPY . .
COPY {{.goRelPath}}/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/{{.exeFile}} {{.goRelPath}}/{{.goFile}}


FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/{{.exeFile}} /app/{{.exeFile}}
COPY --from=builder /app/etc /app/etc

CMD ["./{{.exeFile}}"{{.argument}}]
`
)

func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, map[string]string{
		dockerTemplateFile: dockerTemplate,
	})
}
