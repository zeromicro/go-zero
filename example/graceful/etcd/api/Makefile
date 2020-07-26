version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	docker pull alpine
	docker pull golang:alpine
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/graceful:$(version) . -f example/graceful/etcd/api/Dockerfile
	docker image prune --filter label=stage=gobuilder -f

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/graceful:$(version)

deploy: push
	kubectl -n kevin set image deployment/graceful-deployment graceful=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/graceful:$(version)
