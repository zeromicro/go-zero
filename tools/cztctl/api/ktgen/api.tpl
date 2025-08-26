package {{.Pkg}}

import com.google.gson.Gson

object {{with .Info}}{{.Title}}{{end}}{
	{{range .Types}}
	data class {{.Name}}({{$length := (len .Members)}}{{range $i,$item := .Members}}
		val {{with $item}}{{lowCamelCase .Name}}: {{parseType .Type.Name}}{{end}}{{if ne $i (add $length -1)}},{{end}}{{end}}
	){{end}}
	{{with .Service}}
	{{range .Routes}}suspend fun {{routeToFuncName .Method .Path}}({{with .RequestType}}{{if ne .Name ""}}
		req:{{.Name}},{{end}}{{end}}
		onOk: (({{with .ResponseType}}{{.Name}}{{end}}) -> Unit)? = null,
        onFail: ((String) -> Unit)? = null,
        eventually: (() -> Unit)? = null
    ){
        apiRequest("{{upperCase .Method}}","{{.Path}}",{{with .RequestType}}{{if ne .Name ""}}body=req,{{end}}{{end}} onOk = { {{with .ResponseType}}
            onOk?.invoke({{if ne .Name ""}}Gson().fromJson(it,{{.Name}}::class.java){{end}}){{end}}
        }, onFail = onFail, eventually =eventually)
    }
	{{end}}{{end}}
}
