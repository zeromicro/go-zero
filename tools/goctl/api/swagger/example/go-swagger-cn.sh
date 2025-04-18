#!/bin/bash

# 1. 检查并安装 swagger
if ! command -v swagger &> /dev/null; then
    echo "swagger 未安装，正在从 GitHub 安装..."
    # 这里使用 go-swagger 的安装方式
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest
    if [ $? -ne 0 ]; then
        echo "安装 swagger 失败"
        exit 1
    fi
    echo "swagger 安装成功"
else
    echo "swagger 已安装"
fi

mkdir bin output

export GOBIN=$(pwd)/bin

# 2. 安装最新版 goctl
go install ../../..
if [ $? -ne 0 ]; then
    echo "安装 goctl 失败"
    exit 1
fi
echo "goctl 安装成功"

# 3. 生成 swagger 文件
echo "正在生成 swagger 文件..."
./bin/goctl api swagger --api example_cn.api --dir output
if [ $? -ne 0 ]; then
    echo "生成 swagger 文件失败"
    exit 1
fi

# 4. 启动 swagger 服务
echo "启动 swagger 服务..."
swagger serve ./output/example_cn.json