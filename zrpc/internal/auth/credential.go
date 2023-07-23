package auth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// A Credential is used to authenticate.
type Credential struct {
	App   string
	Token string
}

// GetRequestMetadata gets the request metadata.
func (c *Credential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		appKey:   c.App,
		tokenKey: c.Token,
	}, nil
}

// RequireTransportSecurity always returns false.
func (c *Credential) RequireTransportSecurity() bool {
	return false
}

// ParseCredential parses credential from given ctx.
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
