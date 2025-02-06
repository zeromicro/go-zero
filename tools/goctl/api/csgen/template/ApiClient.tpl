namespace {{.Namespace}};

public sealed class {{.ClientName}}Client : ApiBaseClient
{
    public {{.ClientName}}Client(string host, short port, string scheme = "http") : base(host, port, scheme){}
    {{range $i, $r := .Routes}}
    public async Task<{{if $r.ResponseType}}{{$r.ResponseType}}{{else}}HttpResponseMessage{{end}}> {{$r.HttpMethod}}{{$r.ActionPrefix}}{{$r.ActionName}}Async({{if $r.RequestType}}{{$r.RequestType}} request, {{end}}CancellationToken cancellationToken, HttpContent? body=null)
    {
        return await {{if $r.RequestType}}{{if $r.ResponseType}}RequestResultAsync<{{$r.RequestType}},{{$r.ResponseType}}>(
            HttpMethod.{{$r.HttpMethod}},
            "{{$r.Prefix}}{{$r.UrlPath}}",
            request,{{else}}RequestAsync(
            HttpMethod.{{$r.HttpMethod}},
            "{{$r.Prefix}}{{$r.UrlPath}}",
            request,{{end}}{{else}}{{if $r.ResponseType}}CallResultAsync<{{$r.ResponseType}}>(
            HttpMethod.{{$r.HttpMethod}},
            "{{$r.Prefix}}{{$r.UrlPath}}",{{else}}CallAsync(
            HttpMethod.{{$r.HttpMethod}},
            "{{$r.Prefix}}{{$r.UrlPath}}",{{end}}{{end}}
            cancellationToken,
            body
        );
    }{{end}}
}