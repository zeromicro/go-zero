#!/bin/bash

source build.env
APP_NAME=$APP_NAME
VERSION=$APP_VERSION
BUILD_DIR="dist"
ZIP_DIR="${BUILD_DIR}/zips"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

rm -rf "${BUILD_DIR}"
mkdir -p "${ZIP_DIR}"

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    OUTPUT="${BUILD_DIR}/${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"

    if [ "${GOOS}" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for ${GOOS}/${GOARCH}..."

    env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${OUTPUT}" goctl.go

    if [ $? -ne 0 ]; then
        echo "Error building for ${GOOS}/${GOARCH}"
        exit 1
    fi

    ZIP_OUTPUT="${ZIP_DIR}/$(basename "${OUTPUT}")"
    if [ "${GOOS}" = "windows" ]; then
        zip -j "${ZIP_OUTPUT%.exe}.zip" "${OUTPUT}"
    else
        zip -j "${ZIP_OUTPUT}.zip" "${OUTPUT}"
    fi

    echo "Created zip: ${ZIP_OUTPUT}.zip"
done

echo "All builds completed successfully. Zip files are in ${ZIP_DIR}/"