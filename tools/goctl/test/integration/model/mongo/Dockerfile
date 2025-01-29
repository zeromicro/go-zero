FROM golang:1.22-alpine

# set timezone if you need
# ENV TZ Asia/Shanghai
# if you are in China, you can use the following command to speed up the download
# ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app
COPY goctl /usr/bin/
COPY cmd.sh .

RUN chmod +x /usr/bin/goctl cmd.sh
CMD ["/bin/bash", "cmd.sh"]
