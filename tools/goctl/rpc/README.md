# Rpc Generation
Goctl Rpc是`goctl`脚手架下的一个rpc服务代码生成模块，支持proto模板生成和rpc服务代码生成，通过此工具生成代码你只需要关注业务逻辑编写而不用去编写一些重复性的代码。这使得我们把精力重心放在业务上，从而加快了开发效率且降低了代码出错率。

# 特性
* 简单易用
* 快速提升开发效率
* 出错率低

# 快速开始

### 生成proto模板

```shell script
$ goctl rpc template -o=user.proto
```

```golang
syntax = "proto3";

package remoteuser;

message Request {
  string username = 1;
  string password = 2;
}

message Response {
  string name = 1;
  string gender = 2;
}

service User{
  rpc Login(Request)returns(Response);
}
```
### 生成rpc服务代码

```
$ goctl rpc proto -src=user.proto
```

代码tree

```
```