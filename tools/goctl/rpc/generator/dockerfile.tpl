FROM {{.imageTag}} as builder

# Define the project name | 定义项目名称
ARG PROJECT={{.serviceName}}

WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -ldflags="-s -w" -o /app/${PROJECT}_rpc ${PROJECT}.go

FROM alpine:latest

# Define the project name | 定义项目名称
ARG PROJECT={{.serviceName}}
# Define the config file name | 定义配置文件名
ARG CONFIG_FILE={{.serviceName}}.yaml
# Define the author | 定义作者
ARG AUTHOR=RyanSU@yuansu.china.work@gmail.com

LABEL MAINTAINER=${AUTHOR}

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}

COPY --from=builder /app/${PROJECT}_rpc ./
COPY --from=builder /app/etc/${CONFIG_FILE} ./etc/

EXPOSE 9100
ENTRYPOINT ./${PROJECT}_rpc -f etc/${CONFIG_FILE}