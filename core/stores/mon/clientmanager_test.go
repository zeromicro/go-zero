package mon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func init() {
	_ = mtest.Setup()
}

func TestClientManger_getClient(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		Inject(mtest.ClusterURI(), mt.Client)
		cli, err := getClient(mtest.ClusterURI())
		assert.Nil(t, err)
		assert.Equal(t, mt.Client, cli)
	})
}
