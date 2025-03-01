import {{.Caller}} from "./gocliRequest";{{if .ComponentName}}
import * as components from "./{{.ComponentName}}";
export * from "./{{.ComponentName}}";{{end}}

{{range $i, $r := .Routes}}
{{$r.Comment}}
export function {{$r.FuncName}}{{$r.GenericsTypes}}({{$r.FuncArgs}}) {
    return {{$.Caller}}.{{$r.HttpMethod}}<{{$r.ResponseType}}>({{$r.CallArgs}});
}
{{end}}
