import { ApiBaseClient } from './ApiBaseClient';
{{range $i, $v := .RequestTypes}}
import type { {{$v}} } from './{{$v}}';
{{end}}
{{range $i, $v := .ResponseTypes}}
import type { {{$v}}{{range $i, $s := index $.ResponseSubTypes $v}}, {{$s}}{{end}}} from './{{$v}}';
{{end}}

export class {{.ClientName}}Client extends ApiBaseClient
{
    {{range $i, $r := .Routes}}
    async {{$r.HttpMethod}}{{$r.ActionPrefix}}{{$r.ActionName}}({{if $r.RequestType}}request: {{$r.RequestType}},{{end}}body?: any): Promise<{{if $r.ResponseType}}{{$r.ResponseType}}{{else}}UniApp.RequestSuccessCallbackResult{{end}}> {
        const response = await this.request(
            '{{ToUpper $r.HttpMethod}}',
            '{{$r.Prefix}}{{$r.UrlPath}}',
            {{if $r.RequestHasQueryString}}request.query{{else}}undefined{{end}},
            {{if $r.RequestHasHeaders}}request.header{{else}}undefined{{end}},
            body{{if $r.RequestHasBody}} ?? request.body{{end}}
        );
        {{if $r.ResponseType}}
        const result: {{$r.ResponseType}} = {
            {{if $r.ResponseBodyType}}body: response.data as {{$r.ResponseBodyType}},{{end}}
            {{if $r.ResponseHeadersType}}header: response.header as {{$r.ResponseHeadersType}},{{end}}
        };
        return result;
        {{else}}
        return response;
        {{end}}
    }
    {{end}}
}