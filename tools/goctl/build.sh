git pull

go mod tidy

go build -ldflags="-s -w" -o goctls goctl.go

mv goctls "$GOPATH"/bin/goctls