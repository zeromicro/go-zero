FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update
RUN apk add upx

WORKDIR $GOPATH/src/zero
COPY . .
RUN go build -ldflags="-s -w" -o /app/graceful example/graceful/etcd/api/graceful.go
RUN upx /app/graceful


FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/graceful /app/graceful
COPY example/graceful/etcd/api/etc/graceful-api.json /app/etc/config.json

CMD ["./graceful", "-f", "etc/config.json"]
