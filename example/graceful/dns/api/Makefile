version := v$(shell /bin/date "+%y%m%d%H%M%S")

build:
	docker pull alpine
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/graceful:$(version) . -f example/graceful/dns/api/Dockerfile

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/graceful:$(version)

deploy: push
	kubectl -n kevin set image deployment/graceful-deployment graceful=registry-vpc.cn-hangzhou.aliyuncs.com/xapp/graceful:$(version)
