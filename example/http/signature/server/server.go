package main

import (
	"flag"
	"io"
	"net/http"

	"zero/core/logx"
	"zero/core/service"
	"zero/rest"
	"zero/rest/httpx"
)

var keyPem = flag.String("prikey", "private.pem", "the private key file")

type Request struct {
	User string `form:"user,optional"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := httpx.Parse(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	io.Copy(w, r.Body)
}

func main() {
	flag.Parse()

	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Path: "logs",
			},
		},
		Port: 3333,
		Signature: rest.SignatureConf{
			Strict: true,
			PrivateKeys: []rest.PrivateKeyConf{
				{
					Fingerprint: "bvw8YlnSqb+PoMf3MBbLdQ==",
					KeyFile:     *keyPem,
				},
			},
		},
	})
	defer engine.Stop()

	engine.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/a/b",
		Handler: handle,
	})
	engine.Start()
}
