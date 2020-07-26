FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR $GOPATH/src/zero
COPY . .
RUN go build -ldflags="-s -w" -o /app/main example/load/simulate/cpu/main.go


FROM alpine

RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

RUN apk add git
RUN go get github.com/vikyd/go-cpu-load

RUN mkdir /app
COPY --from=builder /app/main /app/main

WORKDIR /app
CMD ["/app/main"]
