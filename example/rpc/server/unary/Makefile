version := v1

build:
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/unaryserver:$(version) . -f example/rpc/server/unary/Dockerfile
	docker image prune --filter label=stage=gobuilder -f

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/unaryserver:$(version)

deploy: push
	kubectl -n adhoc set image deployment/unaryserver-deployment unaryserver=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/unaryserver:$(version)
