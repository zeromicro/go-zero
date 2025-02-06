{{range $i, $n := .ImportTypes}}
import { {{$n}} } from './{{$n}}';
{{end}}
{{range $i, $s := .SubMessages}}
export type {{$s.MessageName}} = {
    {{range $i, $f := $s.Fields}}{{$f.FieldName}}{{if $f.IsOptional}}?{{end}}: {{$f.TypeName}};{{end}}
};
{{end}}
export type {{.MessageName}} = {
    {{range $i, $f := .Fields}}{{$f.FieldName}}{{if $f.IsOptional}}?{{end}}: {{$f.TypeName}};{{end}}
};