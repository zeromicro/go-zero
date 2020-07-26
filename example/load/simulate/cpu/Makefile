version := v1

build:
	cd $(GOPATH)/src/zero && docker build -t registry.cn-hangzhou.aliyuncs.com/xapp/shedding:$(version) . -f example/load/simulate/cpu/Dockerfile

push: build
	docker push registry.cn-hangzhou.aliyuncs.com/xapp/shedding:$(version)

deploy: push
	kubectl apply -f shedding.yaml

clean:
	kubectl delete -f shedding.yaml
