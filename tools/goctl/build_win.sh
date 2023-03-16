# use git bash

git pull

go mod tidy

GOOS=windows go build -ldflags="-s -w" -o goctls.exe goctl.go

mv goctls.exe $GOPATH/bin/goctls.exe