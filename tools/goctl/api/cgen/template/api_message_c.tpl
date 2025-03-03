#include <stdlib.h>
#include <cJSON.h>
#include "base.h"
#include "message.h"

// headers
{{range $i, $m := .Messages}}{{if $m.HeaderCount}}bool {{$m.MessageName}}_to_curl_slist(curl_slist_t **headers, {{$m.MessageName}}_t *message) {
    char *header = malloc(1024);
    if (NULL == header) {
        return false;
    }
    {{range $j, $f := index $m.Fields "header"}}
    if (message->{{$f.FieldName}}) {
        sprintf(header, "{{$f.FieldTagName}}: {{$f.FieldFormatTag}}", message->{{$f.FieldName}});
        *headers = curl_slist_append(*headers, header);
    }
    {{end}}
    free(header);
    return true;
}
{{end}}{{end}}

// message to json
{{range $i, $m := .Messages}}{{if $m.BodyCount}}char* {{$m.MessageName}}_to_json({{$m.MessageName}}_t *message) {
    char* text = NULL;
{{template "to_object" $m.CJson}}

    text = cJSON_Print(v_{{$m.CJson.VarName}});
end:
    cJSON_Delete(v_{{$m.CJson.VarName}});
    return text;
}
{{end}}{{end}}

// message from json
{{range $i, $m := .Messages}}{{if $m.BodyCount}}bool {{$m.MessageName}}_from_json({{$m.MessageName}}_t *message, char* body) {
    cJSON* v_{{$m.CJson.VarName}} = cJSON_Parse(body);
{{template "from_object" $m.CJson}}

exit_free:
    cJSON_free(v_{{$m.CJson.VarName}});
    return true;
}
{{end}}{{end}}