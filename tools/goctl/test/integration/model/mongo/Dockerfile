FROM golang:1.18

ENV TZ Asia/Shanghai
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app
ADD goctl /usr/bin/goctl
ADD cmd.sh .

RUN chmod +x /usr/bin/goctl
RUN chmod +x cmd.sh
CMD ["/bin/bash", "cmd.sh"]
