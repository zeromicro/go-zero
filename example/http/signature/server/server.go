package main

import (
	"flag"
	"io"
	"net/http"

	"zero/core/logx"
	"zero/core/service"
	"zero/ngin"
	"zero/ngin/httpx"
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

	engine := ngin.MustNewEngine(ngin.NgConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Path: "logs",
			},
		},
		Port: 3333,
		Signature: ngin.SignatureConf{
			Strict: true,
			PrivateKeys: []ngin.PrivateKeyConf{
				{
					Fingerprint: "bvw8YlnSqb+PoMf3MBbLdQ==",
					KeyFile:     *keyPem,
				},
			},
		},
	})
	defer engine.Stop()

	engine.AddRoute(ngin.Route{
		Method:  http.MethodPost,
		Path:    "/a/b",
		Handler: handle,
	})
	engine.Start()
}
