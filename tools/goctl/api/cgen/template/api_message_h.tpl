#ifndef GO_ZERO_API_MESSAGE_H
#define GO_ZERO_API_MESSAGE_H

#include <stdbool.h>
#include <stdint.h>

// declare
{{range $i, $m := .Messages}}{{if $m.FieldCount}}struct __{{$m.MessageName}}_t;
{{end}}{{end}}

// define
{{range $i, $m := .Messages}}{{if $m.FieldCount}}typedef struct __{{$m.MessageName}}_t { {{range $tag, $fs := $m.Fields}}{{range $j, $f := $fs}}
    {{$f.FieldType}} {{$f.FieldName}};{{end}}{{end}}
} {{$m.MessageName}}_t;
{{end}}
{{end}}

// function declare
{{range $i, $m := .Messages}}{{if $m.HeaderCount}}bool {{$m.MessageName}}_to_curl_slist(curl_slist_t **headers, {{$m.MessageName}}_t *message);
{{end}}{{end}}
{{range $i, $m := .Messages}}{{if $m.BodyCount}}char* {{$m.MessageName}}_to_json({{$m.MessageName}}_t *message);
{{end}}{{end}}
{{range $i, $m := .Messages}}{{if $m.BodyCount}}bool {{$m.MessageName}}_from_json({{$m.MessageName}}_t *message, char* body);
{{end}}{{end}}
#endif
