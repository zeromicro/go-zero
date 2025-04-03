#!/bin/bash

is_docker_running() {
    if ! docker info >/dev/null 2>&1; then
        return 1  # Docker is not running
    else
        return 0  # Docker is running
    fi
}

# 1. Check and install Docker if not exists
if ! command -v docker &> /dev/null; then
    echo "Docker not found, attempting to install..."

    # Install Docker using official installation script
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh

    # Verify installation
    if ! command -v docker &> /dev/null; then
        echo "Failed to install Docker"
        exit 1
    fi

    # Add current user to docker group (may require logout/login)
    sudo usermod -aG docker $USER
    echo "Docker installed successfully. You may need to logout and login again for changes to take effect."
else
    echo "Docker already installed"
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

if ! is_docker_running; then
    echo "Docker is not running, Pls start Docker first"
fi

# 4. Clean up any existing swagger-ui container
echo "Cleaning up existing swagger-ui containers..."
docker rm -f swagger-ui 2>/dev/null && echo "Removed existing swagger-ui container"

# 5. Run swagger-ui in Docker
echo "Starting swagger-ui in Docker..."
docker run -d --name swagger-ui -p 8080:8080 -e SWAGGER_JSON=/tmp/example.json -v $(pwd)/example.json:/tmp/example.json swaggerapi/swagger-ui
if [ $? -ne 0 ]; then
    echo "Failed to start swagger-ui container"
    exit 1
fi

echo "Waiting for swagger-ui to initialize..."
sleep 1
SWAGGER_URL="http://localhost:8080"
echo -e "\nSwagger UI is ready at: \033[1;34m${SWAGGER_URL}\033[0m"
echo "Opening in default browser..."

case "$(uname -s)" in
    Linux*)  xdg-open "$SWAGGER_URL";;
    Darwin*) open "$SWAGGER_URL";;
    CYGWIN*|MINGW*|MSYS*) start "$SWAGGER_URL";;
    *) echo "System not supported";;
esac