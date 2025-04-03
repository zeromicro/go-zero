#!/bin/bash

# 检查Docker是否运行的函数
is_docker_running() {
    if ! docker info >/dev/null 2>&1; then
        return 1  # Docker未运行
    else
        return 0  # Docker正在运行
    fi
}

# 1. 检查并安装Docker（如果不存在）
if ! command -v docker &> /dev/null; then
    echo "未检测到Docker，正在尝试安装..."

    # 使用官方脚本安装Docker
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh

    # 验证安装
    if ! command -v docker &> /dev/null; then
        echo "Docker安装失败"
        exit 1
    fi

    # 将当前用户加入docker组（可能需要重新登录）
    sudo usermod -aG docker $USER
    echo "Docker安装成功。您可能需要注销并重新登录使更改生效。"
else
    echo "Docker已安装"
fi

# 2. 检查并安装/验证goctl
if ! command -v goctl &> /dev/null; then
    echo "未检测到goctl，正在从GitHub安装..."
    # 安装最新版goctl
    GOPROXY=https://goproxy.cn,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
    if [ $? -ne 0 ]; then
        echo "goctl安装失败"
        exit 1
    fi
    echo "goctl安装成功"
else
    echo "检测到goctl，正在检查版本..."
    # 获取goctl版本并比较
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
        echo "goctl版本 $version 满足要求 (>= $required)"
    else
        echo "goctl版本 $version 低于要求版本 (>= $required)，正在更新..."
        GOPROXY=https://goproxy.cn,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
        if [ $? -ne 0 ]; then
            echo "goctl更新失败"
            exit 1
        fi
        echo "goctl更新成功"
    fi
fi

# 3. 生成swagger文件
echo "正在生成swagger文件..."
goctl api swagger --api example_cn.api --dir .
if [ $? -ne 0 ]; then
    echo "swagger文件生成失败"
    exit 1
fi

# 检查Docker是否运行
if ! is_docker_running; then
    echo "Docker未运行，请先启动Docker服务"
    exit 1
fi

# 4. 清理现有的swagger-ui容器
echo "正在清理现有的swagger-ui容器..."
docker rm -f swagger-ui 2>/dev/null && echo "已移除现有的swagger-ui容器"

# 5. 在Docker中运行swagger-ui
echo "正在启动swagger-ui容器..."
docker run -d --name swagger-ui -p 8080:8080 \
    -e SWAGGER_JSON=/tmp/example.json \
    -v $(pwd)/example_cn.json:/tmp/example.json \
    swaggerapi/swagger-ui

if [ $? -ne 0 ]; then
    echo "swagger-ui容器启动失败"
    exit 1
fi

# 等待1秒确保服务就绪
echo "等待swagger-ui初始化..."
sleep 1

# 显示访问信息并尝试打开浏览器
SWAGGER_URL="http://localhost:8080"
echo -e "\nSwagger UI 已准备就绪，访问地址: \033[1;34m${SWAGGER_URL}\033[0m"
echo "正在尝试在默认浏览器中打开..."

# 跨平台打开浏览器
case "$(uname -s)" in
    Linux*)  xdg-open "$SWAGGER_URL";;
    Darwin*) open "$SWAGGER_URL";;
    CYGWIN*|MINGW*|MSYS*) start "$SWAGGER_URL";;
    *) echo "无法在当前操作系统自动打开浏览器";;
esac
