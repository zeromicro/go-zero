import httpx
from .base import *
from .message import *

class {{$.ClientName}}:
    def __init__(self, host, port=80, scheme='http'):
        self.__base_url = f'{scheme}://{host}:{port}'

    {{range $i, $a := $.Actions}}
    def {{$a.HttpMethod}}_{{$a.ActionName}}(self{{if $a.RequestMessage}}, req: {{$a.RequestMessage.MessageName}}{{end}}{{if or (not $a.RequestMessage) (not $a.RequestMessage.BodyCount)}}, body=None{{end}}) -> {{if $a.ResponseMessage}}{{$a.ResponseMessage.MessageName}}{{else}}httpx.Response{{end}}:
        r = httpx.{{$a.HttpMethod}}({{if and $a.RequestMessage $a.RequestMessage.PathCount}}req.get_path(f'{self.__base_url}{{$a.UrlPrefix}}{{$a.UrlPath}}'){{else}}f'{self.__base_url}{{$a.UrlPrefix}}{{$a.UrlPath}}'{{end}}{{if $a.RequestMessage }}{{if $a.RequestMessage.HeaderCount}},
            headers=req.get_headers(){{end}}{{if $a.RequestMessage.BodyCount}},
            json=req.get_body(){{else}},
            json=body{{end}}{{if $a.RequestMessage.FormCount}},
            params=req.get_query_params(){{end}}
        {{else}},json=body{{end}})
        if r.status_code != httpx.codes.OK:
            raise OpenApiException(r.content, r.status_code)
        {{if $a.ResponseMessage}}
        result = {{$a.ResponseMessage.MessageName}}(){{if $a.ResponseMessage.HeaderCount}}
        result.set_headers(r.headers){{end}}{{if $a.ResponseMessage.BodyCount}}
        result.set_body(r.json()){{end}}
        return result
        {{else}}
        return r
        {{end}}
    {{end}}