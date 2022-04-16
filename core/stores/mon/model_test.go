package mon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestModel_StartSession(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		m := &Model{
			Collection: newCollection(mt.Coll, breaker.NewBreaker()),
			cli:        mtest.GlobalClient(),
			brk:        breaker.NewBreaker(),
		}
		assert.Panics(t, func() {
			m.StartSession()
		})
	})
}
