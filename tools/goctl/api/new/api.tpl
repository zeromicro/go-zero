type Request {
  Name string `path:"name,options=you|me"`
}

type Response {
  Message string `json:"message"`
}

service {{.name}}-api {
  @handler {{.handler}}Handler
  get /from/:name(Request) returns (Response)
}
