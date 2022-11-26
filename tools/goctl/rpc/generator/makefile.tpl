RED  =  "\e[31;1m"
GREEN = "\e[32;1m"
YELLOW = "\e[33;1m"
BLUE  = "\e[34;1m"
PURPLE = "\e[35;1m"
CYAN  = "\e[36;1m"

docker:
	docker build -f Dockerfile -t ${DOCKER_USERNAME}/{{.serviceName}}-rpc:${VERSION} .
	@printf $(GREEN)"[SUCCESS] build docker successfully"

publish-docker:
	echo "${DOCKER_PASSWORD}" | docker login --username ${DOCKER_USERNAME} --password-stdin http://${REPO}
	docker push ${REPO}/${DOCKER_USERNAME}/{{.serviceName}}-rpc:${VERSION}
	@printf $(GREEN)"[SUCCESS] publish docker successfully"

gen-rpc:
	goctls rpc protoc ./{{.serviceName}}.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	@printf $(GREEN)"[SUCCESS] generate rpc successfully"
{{if .isEnt}}
gen-ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./ent/template/*.tmpl" ./ent/schema
	@printf $(GREEN)"[SUCCESS] generate ent successfully"

gen-rpc-ent-logic:
	goctls rpc ent --schema=./ent/schema  --style=go_zero --multiple=false --serviceName={{.serviceName}} --o=./ --model=$(model)
	@printf $(GREEN)"[SUCCESS] generate ent logic codes successfully"
{{end}}
