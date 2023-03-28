FROM golang:1.20.2-alpine3.17 as builder

WORKDIR /home
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -ldflags="-s -w" -o /home/{{.serviceName}}_api {{.serviceName}}.go

FROM alpine:latest

WORKDIR /home

COPY --from=builder /home/{{.serviceName}}_api ./
COPY --from=builder /home/etc/{{.serviceName}}.yaml ./

EXPOSE {{.port}}
ENTRYPOINT ./{{.serviceName}}_api -f {{.serviceName}}.yaml