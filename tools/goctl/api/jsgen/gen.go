package jsgen

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/tal-tech/go-zero/tools/goctl/api/ktgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const (
	baseTemplate = `
var server='http://localhost:8888'
function apiRequest(method,uri,body,onOk,onFail,eventually){
    var xhr=new XMLHttpRequest();
    xhr.onreadystatechange=function(e){
        if(xhr.readyState==4){
            if(xhr.status==200){
                if(onOk){
                    onOk(JSON.parse(xhr.responseText));
                }
            }else {
                if(onFail){
                    try{
                        onFail(JSON.parse(xhr.responseText))
                    }catch(e){
                        onFail(xhr.responseText)
                    }
                }
            }
            if(eventually){
                eventually()
            }
        }
    }
    xhr.open(method,server+uri)
    xhr.setRequestHeader('Content-Type','application/json')
    xhr.setRequestHeader('Cookies',document.cookie)
    if(body){
        if (typeof body == 'string'){
            xhr.send(body)
        }else{
            xhr.send(JSON.stringify(body))
        }
    }else{
        xhr.send()
    }
}`
	apiTemplate = `{{with .Service}}{{range .Routes}}
function {{routeToFuncName .Method .Path}}({{with .RequestType}}{{if ne .Name ""}}req,{{end}}{{end}}onOk,onFail,eventually){
    apiRequest('{{upperCase .Method}}','{{.Path}}',{{with .RequestType}}{{if ne .Name ""}}req,{{end}}{{end}}onOk,onFail,eventually)
}{{end}}{{end}}`
)

func genBase(dir string, api *spec.ApiSpec) error {
	e := os.MkdirAll(dir, 0755)
	if e != nil {
		return e
	}
	path := filepath.Join(dir, "base_api.js")
	if _, e := os.Stat(path); e == nil {
		fmt.Println("base_api.js already exists , skipped it.")
		return nil
	}

	file, e := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		return e
	}
	defer file.Close()

	_, e = file.WriteString(baseTemplate)
	return e
}

func genApi(dir string, api *spec.ApiSpec) error {
	name := strcase.ToSnake(api.Info.Title + "_api")
	path := filepath.Join(dir, name+".js")
	api.Info.Title = name

	e := os.MkdirAll(dir, 0755)
	if e != nil {
		return e
	}

	file, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		return e
	}
	defer file.Close()

	t, e := template.New("api").Funcs(ktgen.FuncsMap).Parse(apiTemplate)
	if e != nil {
		return e
	}
	return t.Execute(file, api)
}
