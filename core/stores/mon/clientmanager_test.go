package mon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

func TestClientManger_getClient(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	opt := mopt.Client().SetAppName("testAppName")

	mt.Run("test", func(mt *mtest.T) {
		Inject(mtest.ClusterURI(), mt.Client)
		cli, err := getClient(mtest.ClusterURI(), opt)
		assert.Nil(t, err)
		assert.Equal(t, mt.Client, cli)
	})
}
