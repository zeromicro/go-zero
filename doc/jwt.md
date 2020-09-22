# 基于go-zero实现JWT认证

关于JWT是什么，大家可以看看[官网](https://jwt.io/)，一句话介绍下：是可以实现服务器无状态的鉴权认证方案，也是目前最流行的跨域认证解决方案。

要实现JWT认证，我们需要分成如下两个步骤

* 客户端获取JWT token。
* 服务器对客户端带来的JWT token认证。

## 1.  客户端获取JWT Token

我们定义一个协议供客户端调用获取JWT token，我们新建一个目录jwt然后在目录中执行 `goctl api -o jwt.api`，将生成的jwt.api改成如下：

````go
type JwtTokenRequest struct {
}

type JwtTokenResponse struct {
  AccessToken  string `json:"access_token"`
  AccessExpire int64  `json:"access_expire"`
  RefreshAfter int64  `json:"refresh_after"` // 建议客户端刷新token的绝对时间
}

type GetUserRequest struct { 
  UserId string `json:"userId"`
}

type GetUserResponse struct {
  Name string `json:"name"`
}

service jwt-api {
  @server(
    handler: JwtHandler
  )
  post /user/token(JwtTokenRequest) returns (JwtTokenResponse)
}

@server(
  jwt: JwtAuth
)
service jwt-api {
  @server(
    handler: GetUserHandler
  )
  post /user/info(GetUserRequest) returns (GetUserResponse)
}
````

在服务jwt目录中执行：`goctl api go -api jwt.api -dir .`
打开jwtlogic.go文件，修改 `func (l *JwtLogic) Jwt(req types.JwtTokenRequest) (*types.JwtTokenResponse, error) {` 方法如下：

```go

func (l *JwtLogic) Jwt(req types.JwtTokenRequest) (*types.JwtTokenResponse, error) {
	var accessExpire = l.svcCtx.Config.JwtAuth.AccessExpire

	now := time.Now().Unix()
	accessToken, err := l.GenToken(now, l.svcCtx.Config.JwtAuth.AccessSecret, nil, accessExpire)
	if err != nil {
		return nil, err
	}

	return &types.JwtTokenResponse{
    AccessToken:  accessToken,
    AccessExpire: now + accessExpire,
    RefreshAfter: now + accessExpire/2,
  }, nil
}

func (l *JwtLogic) GenToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
```

在启动服务之前，我们需要修改etc/jwt-api.yaml文件如下：
```yaml
Name: jwt-api
Host: 0.0.0.0
Port: 8888
JwtAuth:
  AccessSecret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  AccessExpire: 604800
```
启动服务器，然后测试下获取到的token。

```sh
➜ curl --location --request POST '127.0.0.1:8888/user/token'
{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEyNjE0MjksImlhdCI6MTYwMDY1NjYyOX0.6u_hpE_4m5gcI90taJLZtvfekwUmjrbNJ-5saaDGeQc","access_expire":1601261429,"refresh_after":1600959029}
```

## 2. 服务器验证JWT token

1. 在api文件中通过`jwt: JwtAuth`标记的service表示激活了jwt认证。
2. 可以阅读rest/handler/authhandler.go文件了解服务器jwt实现。
3. 修改getuserlogic.go如下：

```go
func (l *GetUserLogic) GetUser(req types.GetUserRequest) (*types.GetUserResponse, error) {
	return &types.GetUserResponse{Name: "kim"}, nil
}
```

* 我们先不带JWT Authorization header请求头测试下，返回http status code是401，符合预期。

```sh
➜ curl -w  "\nhttp: %{http_code} \n" --location --request POST '127.0.0.1:8888/user/info' \
--header 'Content-Type: application/json' \
--data-raw '{
    "userId": "a"
}'

http: 401
```

* 加上Authorization header请求头测试。

```sh
➜ curl -w  "\nhttp: %{http_code} \n" --location --request POST '127.0.0.1:8888/user/info' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEyNjE0MjksImlhdCI6MTYwMDY1NjYyOX0.6u_hpE_4m5gcI90taJLZtvfekwUmjrbNJ-5saaDGeQc' \
--header 'Content-Type: application/json' \
--data-raw '{
    "userId": "a"
}'
{"name":"kim"}
http: 200
```

综上所述：基于go-zero的JWT认证完成，在真实生产环境部署时候，AccessSecret, AccessExpire, RefreshAfter根据业务场景通过配置文件配置，RefreshAfter 是告诉客户端什么时候该刷新JWT token了，一般都需要设置过期时间前几天。

