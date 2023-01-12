package consul

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/mapping"
)

type target struct {
	Addr              string        `key:",optional"`
	User              string        `key:",optional"`
	Password          string        `key:",optional"`
	Service           string        `key:",optional"`
	Wait              time.Duration `key:"wait,optional"`
	Timeout           time.Duration `key:"timeout,optional"`
	MaxBackoff        time.Duration `key:"max-backoff,optional"`
	Tag               string        `key:"tag,optional"`
	Near              string        `key:"near,optional"`
	Limit             int           `key:"limit,optional"`
	Healthy           bool          `key:"healthy,optional"`
	TLSInsecure       bool          `key:"insecure,optional"`
	Token             string        `key:"token,optional"`
	Dc                string        `key:"dc,optional"`
	AllowStale        bool          `key:"allow-stale,optional"`
	RequireConsistent bool          `key:"require-consistent,optional"`
}

func (t *target) String() string {
	return fmt.Sprintf("service='%s' healthy='%t' tag='%s'", t.Service, t.Healthy, t.Tag)
}

// parseURL with parameters
func parseURL(rawURL url.URL) (target, error) {
	if rawURL.Scheme != schemeName ||
		len(rawURL.Host) == 0 || len(strings.TrimLeft(rawURL.Path, "/")) == 0 {
		return target{},
			errors.Errorf("Malformed URL('%s'). Must be in the next format: 'consul://[user:passwd]@host/service?param=value'", rawURL.String())
	}

	var tgt target
	params := make(map[string]any, len(rawURL.Query()))
	for name, value := range rawURL.Query() {
		params[name] = value[0]
	}
	err := mapping.UnmarshalKey(params, &tgt)
	if err != nil {
		return target{}, errors.Wrap(err, "Malformed URL parameters")
	}

	tgt.User = rawURL.User.Username()
	tgt.Password, _ = rawURL.User.Password()
	tgt.Addr = rawURL.Host
	tgt.Service = strings.TrimLeft(rawURL.Path, "/")

	if len(tgt.Near) == 0 {
		tgt.Near = "_agent"
	}
	if tgt.MaxBackoff == 0 {
		tgt.MaxBackoff = time.Second
	}

	return tgt, nil
}

// consulConfig returns config based on the parsed target.
// It uses custom http-client.
func (t *target) consulConfig() *api.Config {
	var creds *api.HttpBasicAuth
	if len(t.User) > 0 && len(t.Password) > 0 {
		creds = new(api.HttpBasicAuth)
		creds.Password = t.Password
		creds.Username = t.User
	}
	// custom http.Client
	c := &http.Client{
		Timeout: t.Timeout,
	}
	return &api.Config{
		Address:    t.Addr,
		HttpAuth:   creds,
		WaitTime:   t.Wait,
		HttpClient: c,
		TLSConfig: api.TLSConfig{
			InsecureSkipVerify: t.TLSInsecure,
		},
		Token: t.Token,
	}
}
