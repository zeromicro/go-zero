package mon

import (
	"context"
	"io"

	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

var clientManager = syncx.NewResourceManager()

// ClosableClient wraps *mongo.Client and provides a Close method.
type ClosableClient struct {
	*mongo.Client
}

// Close disconnects the underlying *mongo.Client.
func (cs *ClosableClient) Close() error {
	return cs.Client.Disconnect(context.Background())
}

// Inject injects a *mongo.Client into the client manager.
// Typically, this is used to inject a *mongo.Client for test purpose.
func Inject(key string, client *mongo.Client) {
	clientManager.Inject(key, &ClosableClient{client})
}

func getClient(url string, opts ...Option) (*mongo.Client, error) {
	val, err := clientManager.GetResource(url, func() (io.Closer, error) {
		o := mopt.Client().ApplyURI(url)
		opts = append([]Option{defaultTimeoutOption()}, opts...)
		for _, opt := range opts {
			opt(o)
		}

		cli, err := mongo.Connect(context.Background(), o)
		if err != nil {
			return nil, err
		}

		err = cli.Ping(context.Background(), nil)
		if err != nil {
			return nil, err
		}

		concurrentSess := &ClosableClient{
			Client: cli,
		}

		return concurrentSess, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*ClosableClient).Client, nil
}
