package auth

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const defaultExpiration = 5 * time.Minute

// An Authenticator is used to authenticate the rpc requests.
type Authenticator struct {
	store  *redis.Redis
	key    string
	cache  *collection.Cache
	strict bool
}

// NewAuthenticator returns an Authenticator.
func NewAuthenticator(store *redis.Redis, key string, strict bool) (*Authenticator, error) {
	cache, err := collection.NewCache(defaultExpiration)
	if err != nil {
		return nil, err
	}

	return &Authenticator{
		store:  store,
		key:    key,
		cache:  cache,
		strict: strict,
	}, nil
}

// Authenticate authenticates the given ctx.
func (a *Authenticator) Authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, missingMetadata)
	}

	apps, tokens := md[appKey], md[tokenKey]
	if len(apps) == 0 || len(tokens) == 0 {
		return status.Error(codes.Unauthenticated, missingMetadata)
	}

	app, token := apps[0], tokens[0]
	if len(app) == 0 || len(token) == 0 {
		return status.Error(codes.Unauthenticated, missingMetadata)
	}

	return a.validate(app, token)
}

func (a *Authenticator) validate(app, token string) error {
	expect, err := a.cache.Take(app, func() (any, error) {
		return a.store.Hget(a.key, app)
	})
	if err != nil {
		if a.strict {
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	}

	if token != expect {
		return status.Error(codes.Unauthenticated, accessDenied)
	}

	return nil
}
