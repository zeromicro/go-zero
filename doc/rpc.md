# rpc设计规范

* 目录结构
  * service/remote目录下按照服务所属模块存放，比如用户的profile接口，目录如下：
  
     `service/remote/user/profile.proto`
  
  * 生成的profile.pb.go也放在该目录下，并且profile.proto文件里要加上`package user;`

* 错误处理
  * 需要使用status.Error(code, desc)来定义返回的错误
  * code是codes.Code类型，尽可能使用grpc已经定义好的code
  * codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss错误会被自动熔断