#!/bin/bash

# 1. Check and install swagger if not exists
if ! command -v swagger &> /dev/null; then
    echo "swagger not found, installing from GitHub..."
    # Using go-swagger installation method
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest
    if [ $? -ne 0 ]; then
        echo "Failed to install swagger"
        exit 1
    fi
    echo "swagger installed successfully"
else
    echo "swagger already installed"
fi

# 2. Check and install/verify goctl
if ! command -v goctl &> /dev/null; then
    echo "goctl not found, installing from GitHub..."
    # Install latest goctl version
    GOPROXY=https://goproxy.cn,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
    if [ $? -ne 0 ]; then
        echo "Failed to install goctl"
        exit 1
    fi
    echo "goctl installed successfully"
else
    echo "goctl found, checking version..."
    # Get goctl version and compare
    version=$(goctl --version | awk '{print $3}' | tr -d 'v')
    required="1.8.2"

    # Version comparison function
    version_compare() {
        if [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$1" ]; then
            return 0  # $1 >= $2
        else
            return 1  # $1 < $2
        fi
    }

    if version_compare "$version" "$required"; then
        echo "goctl version $version meets requirement (>= $required)"
    else
        echo "goctl version $version is lower than required (>= $required), updating..."
        GOPROXY=https://goproxy.cn,direct go install github.com/zeromicro/go-zero/tools/goctl@latest
        if [ $? -ne 0 ]; then
            echo "Failed to update goctl"
            exit 1
        fi
        echo "goctl updated successfully"
    fi
fi

# 3. Generate swagger files
echo "Generating swagger files..."
goctl api swagger --api example.api --dir .
if [ $? -ne 0 ]; then
    echo "Failed to generate swagger files"
    exit 1
fi

# 4. Start swagger server
echo "Starting swagger server..."
swagger serve example.json