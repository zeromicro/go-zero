go mod tidy

go build -o goctls goctl.go

cp ./goctls $GOPATH/bin/goctls