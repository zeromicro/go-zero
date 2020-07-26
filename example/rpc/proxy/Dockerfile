FROM golang:1.11 AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR $GOPATH/src/zero
COPY . .
RUN go build -ldflags="-s -w" -o /app/unaryproxy example/rpc/proxy/proxy.go


FROM alpine

WORKDIR /app
COPY --from=builder /app/unaryproxy /app/unaryproxy

CMD ["./unaryproxy"]
