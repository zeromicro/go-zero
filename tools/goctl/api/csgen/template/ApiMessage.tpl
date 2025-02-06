using System.Text.Json.Serialization;

namespace {{.Namespace}};

public class {{.MessageName}}
{
    {{range $i, $f := .Fields}}{{if eq $f.Tag "json"}}
    [JsonPropertyName("{{$f.KeyName}}")]
    {{else if eq $f.Tag "header"}}
    [JsonIgnore]
    [HeaderPropertyName("{{$f.KeyName}}")]
    {{else if eq $f.Tag "form"}}
    [JsonIgnore]
    [FormPropertyName("{{$f.KeyName}}")]
    {{else if eq $f.Tag "path"}}
    [JsonIgnore]
    [PathPropertyName("{{$f.KeyName}}")]
    {{end}}
    public {{$f.TypeName}}{{if $f.IsOptional}}?{{end}} {{$f.FieldName}} { get; set; }
    {{end}}
}