FROM golang:1.21-alpine

ENV TZ Asia/Shanghai
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app
COPY goctl /usr/bin/
COPY cmd.sh .

RUN chmod +x /usr/bin/goctl cmd.sh
CMD ["/bin/bash", "cmd.sh"]
