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

# 2. 检查并安装/验证 goctl
if ! command -v goctl &> /dev/null; then
    echo "goctl 未安装，正在从 GitHub 安装..."
    # 安装最新版 goctl
    GOPROXY=https://goproxy.cn,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
    if [ $? -ne 0 ]; then
        echo "安装 goctl 失败"
        exit 1
    fi
    echo "goctl 安装成功"
else
    echo "goctl 已安装，正在检查版本..."
    # 获取 goctl 版本并比较
    version=$(goctl --version | awk '{print $3}' | tr -d 'v')
    required="1.8.2"

    # 版本比较函数
    version_compare() {
        if [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$1" ]; then
            return 0  # $1 >= $2
        else
            return 1  # $1 < $2
        fi
    }

    if version_compare "$version" "$required"; then
        echo "goctl 版本 $version 满足要求 (>= $required)"
    else
        echo "goctl 版本 $version 低于要求 (>= $required)，正在更新..."
        GOPROXY=https://goproxy.cn,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
        if [ $? -ne 0 ]; then
            echo "更新 goctl 失败"
            exit 1
        fi
        echo "goctl 更新成功"
    fi
fi

# 3. 生成 swagger 文件
echo "正在生成 swagger 文件..."
goctl api swagger --api example_cn.api --dir .
if [ $? -ne 0 ]; then
    echo "生成 swagger 文件失败"
    exit 1
fi

# 4. 启动 swagger 服务
echo "启动 swagger 服务..."
swagger serve example_cn.json