RED  =  "\e[31;1m"
GREEN = "\e[32;1m"
YELLOW = "\e[33;1m"
BLUE  = "\e[34;1m"
PURPLE = "\e[35;1m"
CYAN  = "\e[36;1m"

docker:
	docker build -f Dockerfile -t ${DOCKER_USERNAME}/{{.serviceName}}-api:${VERSION} .
	@printf $(GREEN)"[SUCCESS] build docker successfully"

publish-docker:
	echo "${DOCKER_PASSWORD}" | docker login --username ${DOCKER_USERNAME} --password-stdin https://${REPO}
	docker push ${DOCKER_USERNAME}/{{.serviceName}}-api:${VERSION}
	@printf $(GREEN)"[SUCCESS] publish docker successfully"

gen-api:
	goctls api go --api ./desc/all.api --dir ./ --trans_err=true
	swagger generate spec --output=./{{.serviceName}}.yml --scan-models
	@printf $(GREEN)"[SUCCESS] generate API successfully"

gen-swagger:
	swagger generate spec --output=./{{.serviceName}}.yml --scan-models
	@printf $(GREEN)"[SUCCESS] generate swagger successfully"

serve-swagger:
	lsof -i:36666 | awk 'NR!=1 {print $2}' | xargs killall -9 || true
	swagger serve -F=swagger --port 36666 {{.serviceName}}.yml
	@printf $(GREEN)"[SUCCESS] serve swagger-ui successfully"
