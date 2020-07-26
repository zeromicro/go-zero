version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/pub:$(version) . -f example/etcd/pub/Dockerfile
	docker image prune --filter label=stage=gobuilder -f

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/pub:$(version)

deploy: push
	kubectl -n adhoc set image deployment/pub-deployment pub=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/pub:$(version)
