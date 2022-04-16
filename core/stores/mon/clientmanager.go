package mon

import (
	"context"
	"io"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const defaultTimeout = time.Second

var clientManager = syncx.NewResourceManager()

type ClosableClient struct {
	*mongo.Client
}

func (cs *ClosableClient) Close() error {
	return cs.Client.Disconnect(context.Background())
}

func getClient(url string) (*mongo.Client, error) {
	val, err := clientManager.GetResource(url, func() (io.Closer, error) {
		cli, err := mongo.Connect(context.Background(), mopt.Client().ApplyURI(url))
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
