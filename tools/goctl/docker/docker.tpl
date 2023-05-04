FROM {{.Image}} as builder

# Define the project name | 定义项目名称
ARG PROJECT={{.ServiceName}}

WORKDIR /build
COPY . .
{{if .Chinese}}
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
{{end}}{{if .HasTimezone}}
RUN apk update --no-cache && apk add --no-cache tzdata
{{end}}
RUN go env -w GO111MODULE=on \
{{if .Chinese}}    && go env -w GOPROXY=https://goproxy.cn,direct \
{{end}}    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -ldflags="-s -w" -o /build/${PROJECT}_{{.ServiceType}} ${PROJECT}.go

FROM {{.BaseImage}}

# Define the project name | 定义项目名称
ARG PROJECT={{.ServiceName}}
# Define the config file name | 定义配置文件名
ARG CONFIG_FILE={{.ServiceName}}.yaml
# Define the author | 定义作者
ARG AUTHOR="{{.Author}}"

LABEL org.opencontainers.image.authors=${AUTHOR}

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}
{{if .HasTimezone}}
COPY --from=builder /usr/share/zoneinfo/{{.Timezone}} /usr/share/zoneinfo/{{.Timezone}}
ENV TZ={{.Timezone}}
{{end}}
COPY --from=builder /build/${PROJECT}_{{.ServiceType}} ./
COPY --from=builder /build/etc/${CONFIG_FILE} ./etc/
{{if .HasPort}}
EXPOSE {{.Port}}
{{end}}
ENTRYPOINT ./${PROJECT}_{{.ServiceType}} -f etc/${CONFIG_FILE}