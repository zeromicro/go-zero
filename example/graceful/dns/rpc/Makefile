version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	docker pull alpine
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/gracefulrpc:$(version) . -f example/graceful/dns/rpc/Dockerfile

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/gracefulrpc:$(version)

deploy: push
	kubectl -n kevin set image deployment/gracefulrpc-deployment gracefulrpc=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/gracefulrpc:$(version)
