#ifndef GO_ZERO_API_CLIENT_H
#define GO_ZERO_API_CLIENT_H

#include "base.h"
#include "message.h"

// declare
{{range $i, $a := .Actions}}bool {{$.ClientName}}_{{ToLower $a.HttpMethod}}_{{$a.ActionName}}({{$.ClientName}}_t *client{{if $a.RequestMessage}}, {{$a.RequestMessage.MessageName}}_t* req{{end}}{{if $a.ResponseMessage}}, {{$a.ResponseMessage.MessageName}}_t *resp{{end}}{{if not (and $a.RequestMessage $a.RequestMessage.BodyCount)}}, const char *body{{end}});
{{end}}

#endif
