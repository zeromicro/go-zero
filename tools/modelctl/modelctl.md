# modelctl使用说明

## modelctl用途
* 根据数据库中表名生成model.go代码,目前支持通过【指定配置文件】和【命令行参数】两种形式来生成

## modelctl使用说明
* modelctl参数说明

  > 生成配置文件模板  

  `modelctl model template`
  
  > 根据指定配置文件生成*model.go,-c参数为配置文件名称  
  
  参考命令：`modelctl -c config.json`
  
  > 根据命令行生成*model.go  
  
  `modelctl cmd [--address|-a,--schema|-s,--force|-f,--redis|-r]`  
  
  参考命令：`modelctl cmd -a root:123456@127.0.0.1:3306 -s user -f -r `
  
  `--address|-a` 数据库连接地址，格式：[username]:[password]@[address],参考格式：root:123456@127.0.0.1:3306  
  
  `--schema|-s` 指定数据库名称
  
  `--force|-f` 是否强制覆盖源文件，默认：false，强制覆盖将导致原或已修改文件丢失
  
  `--redis|-r` 是否生成redis缓存逻辑代码，默认：false  
   
  详细说明见 `--help|-h`
  
* 配置文件模板说明 

  ```
  {
    "WithCache": false,
    "Force": true,
    "Username": "***",
    "Password": "***",
    "Address": "**",
    "TableSchema":"*",
    "Tables": [
        "**"
    ]
  }
  ```
  `WithCache` 生成文件时是否待redis缓存逻辑代码  
  
  `Force` 是否强制覆盖原有同名文件，覆盖则会丢失原文件  
  
  `Username` 数据库访问用户名  
  
  `Password` 数据库访问用户密码 
  
  `Address` 数据库访问地址  
  
  `TableSchema` 数据库名  
  
  `Tables` 指定生成model的表名，不填或空则按照该库下全部表进行生成  