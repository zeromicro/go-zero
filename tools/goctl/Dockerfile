FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
# if you are in China, you can use the following command to speed up the download
# ENV GOPROXY=https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN addgroup -g 1000 -S app && adduser -u 1000 -S app -G app

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/goctl ./goctl.go


FROM golang:alpine

RUN apk update --no-cache && apk add --no-cache protoc

COPY --from=builder /etc/passwd /etc/group /etc/
COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=builder --chown=1000:1000 /go/bin/protoc-gen-go* /app/goctl /usr/local/bin/
ENV TZ=Asia/Shanghai

WORKDIR /app
USER app

LABEL org.opencontainers.image.authors="Kevin Wan"
LABEL org.opencontainers.image.base.name="docker.io/library/golang:alpine"
LABEL org.opencontainers.image.description="A cloud-native Go microservices framework with cli tool for productivity."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/zeromicro/go-zero"
LABEL org.opencontainers.image.title="goctl (cli)"

ENTRYPOINT ["/usr/local/bin/goctl"]
