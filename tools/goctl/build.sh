go mod tidy

go build -o goctls goctl.go

mv goctls "$GOPATH"/bin/goctls