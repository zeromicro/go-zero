variables:
  VERSION: 0.0.1
  REPO: docker.io

stages:
  - info
  - build
  - publish
  - clean

info-job:
  stage: info
  script:
    - echo "Start build version $VERSION"
    - export VERSION=$VERSION
    - export DOCKER_USERNAME=$DOCKER_USERNAME
    - export DOCKER_PASSWORD=$DOCKER_PASSWORD
    - export REPO=$REPO

build-job:
  stage: build
  script:
    - echo "Compiling the code and build docker image..."
    - make docker
    - echo "Compile complete."

deploy-job:
  stage: publish
  environment: production
  script:
    - echo "Publish docker images..."
    - make publish-docker
    - echo "Docker images successfully published."

clean-job:
  stage: clean
  script:
    # 删除所有 none 镜像 | delete all none images
    - docker images |grep none|awk '{print $3}'|xargs docker rmi
    - echo "Delete all none images successfully."