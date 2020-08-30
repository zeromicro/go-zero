package gen

const dockerTemplate = `FROM golang:alpine AS builder

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
