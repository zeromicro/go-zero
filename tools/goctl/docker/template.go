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
{{if .Chinese}}ENV GOPROXY https://goproxy.cn,direct{{end}}

WORKDIR /build/zero

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY {{.GoRelPath}}/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/{{.ExeFile}} {{.GoRelPath}}/{{.GoFile}}


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/{{.ExeFile}} /app/{{.ExeFile}}
COPY --from=builder /app/etc /app/etc

CMD ["./{{.ExeFile}}"{{.Argument}}]
`
)

func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, map[string]string{
		dockerTemplateFile: dockerTemplate,
	})
}
