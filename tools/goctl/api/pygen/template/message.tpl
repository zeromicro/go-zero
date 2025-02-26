{{range $i, $m := .Messages}}
class {{$m.MessageName}}:
    def __init__(self):{{if $m.BodyCount}}
        self.__json = {}{{end}}{{if $m.HeaderCount}}
        self.__header = {}{{end}}{{if $m.FormCount}}
        self.__form = {}{{end}}{{if $m.PathCount}}
        self.__path = {}{{end}}

    {{range $j, $f := $m.Fields}}@property
    def {{$f.FieldName}}(self):
        return self.__{{$f.FieldTag}}['{{$f.FieldTagName}}']

    @{{$f.FieldName}}.setter
    def {{$f.FieldName}}(self, v):
        self.__{{$f.FieldTag}}['{{$f.FieldTagName}}'] = v
    {{end}}
    {{if $m.BodyCount}}def get_body(self):
        return self.__json
    def set_body(self, v):
        self.__json = v
    {{end}}{{if $m.HeaderCount}}
    def get_headers(self):
        return self.__header
    def set_headers(self, v):
        self.__header = v
    {{end}}{{if $m.FormCount}}
    def get_query_params(self):
        return self.__form
    {{end}}{{if $m.PathCount}}
    def get_path(self, path):
        for k, v in self.__path:
            path = path.replace(':' + k, v)
        return path
{{end}}{{end}}