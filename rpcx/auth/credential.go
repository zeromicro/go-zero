package auth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type Credential struct {
	App   string
	Token string
}

func (c *Credential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		appKey:   c.App,
		tokenKey: c.Token,
	}, nil
}

func (c *Credential) RequireTransportSecurity() bool {
	return false
}

func ParseCredential(ctx context.Context) Credential {
	var credential Credential

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return credential
	}

	apps, tokens := md[appKey], md[tokenKey]
	if len(apps) == 0 || len(tokens) == 0 {
		return credential
	}

	app, token := apps[0], tokens[0]
	if len(app) == 0 || len(token) == 0 {
		return credential
	}

	credential.App = app
	credential.Token = token

	return credential
}
