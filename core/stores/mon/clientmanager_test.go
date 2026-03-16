package mon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestClientManger_getClient(t *testing.T) {
	c := &mongo.Client{}
	Inject("foo", c)
	cli, err := getClient("foo")
	assert.Nil(t, err)
	assert.Equal(t, c, cli)
}
