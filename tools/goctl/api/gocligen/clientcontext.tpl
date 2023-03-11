package {{.pkgName}}

import (
	{{.imports}}
)

type ClientContext struct {
	httpc.Service
	host    string
	errType reflect.Type
}

var (
	ErrHostEmpty   = errors.New("host is empty")
	ErrHostInvalid = errors.New("host is invalid")
)

func NewClientContext(cli httpc.Service, host string) (*ClientContext, error) {
	if len(host) == 0 {
		return nil, ErrHostEmpty
	}
	host, err := checkHost(host)
	if err != nil {
		return nil, ErrHostInvalid
	}
	return &ClientContext{
		Service: cli,
		host:    host,
	}, nil
}

func (cc *ClientContext) SetError(err error) {
	if err == nil {
		return
	}
	cc.errType = reflect.TypeOf(err).Elem()
}

func (cc *ClientContext) SetClient(cli *http.Client) {
	cc.Service = httpc.NewServiceWithClient("{{.service}}", cli)
}

func (cc *ClientContext) Host() string {
	return cc.host
}

func checkHost(host string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		host = "http://" + host
	}
	return host, nil
}
