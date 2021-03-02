# 设计思路
支持通过Kubernetes Service进行服务发现，方便集成现有服务。<br />
每个k8s Service都关联有唯一的Endpoints对象，保存了所有ready和not ready的Pod，
deployment扩缩容时，会实时更新Endpoints下的地址列表，使用k8s的informer sdk 
来watch我们感兴趣的Endpoints，再更新到gRPC。



## 配置方法
复用Endpoints配置项，要求数组长度为且格式如下：
```yaml
Transform:
   Endpoints:
   - k8s:///transform-svc.ns:8081
```