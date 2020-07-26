FROM golang:1.13 AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR $GOPATH/src/zero
COPY . .
RUN go build -ldflags="-s -w" -o /app/gracefulrpc example/graceful/dns/rpc/gracefulrpc.go


FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/gracefulrpc /app/gracefulrpc
COPY example/graceful/dns/rpc/etc/config.json /app/etc/config.json

CMD ["./gracefulrpc", "-f", "etc/config.json"]
