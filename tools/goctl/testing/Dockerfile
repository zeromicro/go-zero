FROM golang:latest

WORKDIR /go/src/app

RUN apt update
RUN apt install -y tree
RUN apt install -y protobuf-compiler
RUN go get -u github.com/golang/protobuf/protoc-gen-go

COPY goctl .
COPY test.api ./test/
COPY test.proto ./test-rpc/

CMD ["/bin/sh", "-c", "while true; do sleep 30; done;"]