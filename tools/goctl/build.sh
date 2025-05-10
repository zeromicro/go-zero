#!/bin/bash

# 设置变量
source build.env
APP_NAME=$APP_NAME
VERSION=$APP_VERSION
BUILD_DIR="dist"
ZIP_DIR="${BUILD_DIR}/zips"

# 支持的平台和架构组合
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# 清理并创建构建目录
rm -rf "${BUILD_DIR}"
mkdir -p "${ZIP_DIR}"

# 遍历所有平台进行构建
for PLATFORM in "${PLATFORMS[@]}"; do
    # 分割平台和架构
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    # 设置输出文件名
    OUTPUT="${BUILD_DIR}/${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"

    # 为Windows添加.exe后缀
    if [ "${GOOS}" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for ${GOOS}/${GOARCH}..."

    # 执行构建
    env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${OUTPUT}" goctl.go

    if [ $? -ne 0 ]; then
        echo "Error building for ${GOOS}/${GOARCH}"
        exit 1
    fi

    # 创建zip包
    ZIP_OUTPUT="${ZIP_DIR}/$(basename "${OUTPUT}")"
    if [ "${GOOS}" = "windows" ]; then
        zip -j "${ZIP_OUTPUT%.exe}.zip" "${OUTPUT}"
    else
        zip -j "${ZIP_OUTPUT}.zip" "${OUTPUT}"
    fi

    echo "Created zip: ${ZIP_OUTPUT}.zip"
done

echo "All builds completed successfully. Zip files are in ${ZIP_DIR}/"