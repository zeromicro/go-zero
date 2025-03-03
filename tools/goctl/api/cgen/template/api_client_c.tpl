#include "client.h"
#include <stdlib.h>
#include <string.h>

// define
{{range $i, $a := .Actions}}
bool {{$.ClientName}}_{{ToLower $a.HttpMethod}}_{{$a.ActionName}}({{$.ClientName}}_t *client{{if $a.RequestMessage}}, {{$a.RequestMessage.MessageName}}_t* req{{end}}{{if $a.ResponseMessage}}, {{$a.ResponseMessage.MessageName}}_t *resp{{end}}{{if not (and $a.RequestMessage $a.RequestMessage.BodyCount)}}, const char *body{{end}}) {
    //request headers
    curl_slist_t * headers = NULL;
    headers = curl_slist_append(headers, "Content-Type: application/json");{{if and $a.RequestMessage $a.RequestMessage.HeaderCount}}
    if (!{{$a.RequestMessage.MessageName}}_to_curl_slist(&headers, req)) {
        return false;
    }{{end}}

    {{if and $a.RequestMessage $a.RequestMessage.BodyCount}}// request body
    char *body = {{$a.RequestMessage.MessageName}}_to_json(req);
    {{end}}
    // request path
    char *path = "{{$a.UrlPrefix}}{{$a.UrlPath}}";{{if and $a.RequestMessage $a.RequestMessage.PathCount}}
    {{range $i, $f := index $a.RequestMessage.Fields "path"}} path = replace_substr(path, ":{{$f.FieldTagName}}", req->{{$f.FieldName}}); {{end}}{{end}}

    // request
    base_request_t request = {
        .method = HTTP_METHOD_{{ToUpper $a.HttpMethod}},
        .headers = headers,
        .body = body,
        .body_size = strlen(body),
        .path = path,
    };
    base_response_t response;
    if (!base_client_request(client, &request, &response)) {
        return false;
    }

    {{if $a.ResponseMessage}}// response body
    if (!{{$a.ResponseMessage.MessageName}}_from_json(resp, response.body)) {
        return false;
    }{{end}}

    base_response_free(&response);

    {{if and $a.RequestMessage $a.RequestMessage.BodyCount}}free(body);{{end}}
    
    return true;
}
{{end}}