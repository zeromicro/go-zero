version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	#docker pull alpine
	#docker pull golang:alpine
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/etcdmon:$(version) . -f example/etcd/demo/Dockerfile
	#docker image prune --filter label=stage=gobuilder -f

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/etcdmon:$(version)

deploy: push
	kubectl -n xx-xiaoheiban set image deployment/etcdmon-deployment etcdmon=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/etcdmon:$(version)
