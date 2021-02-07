# 设计思路
支持通过Kubernetes Service进行服务发现，方便集成现有服务。

## 接口抽象
- 抽象subscriber接口，方便后续从支持更多的注册中心发现服务
- publisher暂不作抽象，新的服务注册建议通过etcd或者k8s service