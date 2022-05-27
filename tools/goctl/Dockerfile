FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/goctl ./goctl.go


FROM golang:alpine

RUN apk update --no-cache && apk add --no-cache protoc

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /go/bin/protoc-gen-go /usr/bin/protoc-gen-go
COPY --from=builder /go/bin/protoc-gen-go-grpc /usr/bin/protoc-gen-go-grpc
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/goctl /usr/bin/goctl

CMD ["goctl"]
