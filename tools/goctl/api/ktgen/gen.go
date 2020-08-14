package ktgen

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"github.com/iancoleman/strcase"
)

const (
	apiBaseTemplate = `package {{.}}

import com.google.gson.Gson
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

const val SERVER = "http://localhost:8080"

suspend fun apiPost(
    uri: String,
    body: Any,
    onOk: ((String) -> Unit)? = null,
    onFail: ((String) -> Unit)? = null,
    eventually: (() -> Unit)? = null
) = withContext(Dispatchers.IO) {
    val url = URL(SERVER + uri)
    with(url.openConnection() as HttpURLConnection) {
        requestMethod = "POST"
        headerFields["Content-Type"] = listOf("Application/json")

        val data = when (body) {
            is String -> {
                body
            }
            else -> {
                Gson().toJson(body)
            }
        }
        val wr = OutputStreamWriter(outputStream)
        wr.write(data)
        wr.flush()

        //response
        BufferedReader(InputStreamReader(inputStream)).use {
            val response = it.readText()
            if (responseCode == 200) {
                onOk?.invoke(response)
            } else {
                onFail?.invoke(response)
            }
        }
    }
    eventually?.invoke()
}

suspend fun apiGet(
    uri: String,
    onOk: ((String) -> Unit)? = null,
    onFail: ((String) -> Unit)? = null,
    eventually: (() -> Unit)? = null
) = withContext(Dispatchers.IO) {
    val url = URL(SERVER + uri)
    with(url.openConnection() as HttpURLConnection) {
        requestMethod = "POST"
        headerFields["Content-Type"] = listOf("Application/json")

        val wr = OutputStreamWriter(outputStream)
        wr.flush()

        //response
        BufferedReader(InputStreamReader(inputStream)).use {
            val response = it.readText()
            if (responseCode == 200) {
                onOk?.invoke(response)
            } else {
                onFail?.invoke(response)
            }
        }
    }
    eventually?.invoke()
}
`
	apiTemplate = `package {{with .Info}}{{.Title}}{{end}}

import com.google.gson.Gson

object Api{
	{{range .Types}}
	data class {{.Name}}({{$length := (len .Members)}}{{range $i,$item := .Members}}
		val {{with $item}}{{lowCamelCase .Name}}: {{parseType .Type}}{{end}}{{if ne $i (add $length -1)}},{{end}}{{end}}
	){{end}}
	{{with .Service}}
	{{range .Routes}}suspend fun {{pathToFuncName .Path}}({{if ne .Method "get"}}
		req:{{with .RequestType}}{{.Name}},{{end}}{{end}}
		onOk: (({{with .ResponseType}}{{.Name}}{{end}}) -> Unit)? = null,
        onFail: ((String) -> Unit)? = null,
        eventually: (() -> Unit)? = null
    ){
        api{{if eq .Method "get"}}Get{{else}}Post{{end}}("{{.Path}}",{{if ne .Method "get"}}req,{{end}} onOk = {
            onOk?.invoke(Gson().fromJson(it,{{with .ResponseType}}{{.Name}}{{end}}::class.java))
        }, onFail = onFail, eventually =eventually)
    }
	{{end}}{{end}}
}`
)

func genBase(dir, pkg string, api *spec.ApiSpec) error {
	e := os.MkdirAll(dir, 0755)
	if e != nil {
		logx.Error(e)
		return e
	}
	path := filepath.Join(dir, "BaseApi.kt")
	if _, e := os.Stat(path); e == nil {
		return nil
	}

	file, e := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		logx.Error(e)
		return e
	}
	defer file.Close()

	t, e := template.New("n").Parse(apiBaseTemplate)
	if e != nil {
		logx.Error(e)
		return e
	}
	e = t.Execute(file, pkg)
	if e != nil {
		logx.Error(e)
		return e
	}
	return nil
}

func genApi(dir, pkg string, api *spec.ApiSpec) error {
	path := filepath.Join(dir, strcase.ToCamel(api.Info.Title+"Api")+".kt")
	api.Info.Title= pkg

	e:=os.MkdirAll(dir,0755)
	if e!=nil {
	    logx.Error(e)
	    return e
	}

	file,e:=os.OpenFile(path,os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0644)
	if e!=nil {
	    logx.Error(e)
	    return e
	}
	defer file.Close()

	t,e:=template.New("api").Funcs(funcsMap).Parse(apiTemplate)
	if e!=nil{
	    log.Fatal(e)
	}
	e=t.Execute(file,api)
	if e!=nil{
	    log.Fatal(e)
	}
	return nil
}
