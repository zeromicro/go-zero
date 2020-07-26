version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/unarydirect:$(version) . -f example/rpc/client/direct/Dockerfile

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/unarydirect:$(version)

deploy: push
	kubectl -n adhoc set image deployment/unarydirect-deployment unarydirect=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/unarydirect:$(version)
