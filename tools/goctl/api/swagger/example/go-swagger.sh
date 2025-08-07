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

mkdir bin output

export GOBIN=$(pwd)/bin

# 2. Install latest goctl version
go install ../../..
if [ $? -ne 0 ]; then
    echo "Failed to install goctl"
    exit 1
fi
echo "goctl installed successfully"

# 3. Generate swagger files
echo "Generating swagger files..."
./bin/goctl api swagger --api example.api --dir output
if [ $? -ne 0 ]; then
    echo "Failed to generate swagger files"
    exit 1
fi

# 4. Start swagger server
echo "Starting swagger server..."
swagger serve ./output/example.json