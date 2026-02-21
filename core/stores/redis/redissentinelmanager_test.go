package redis

import (
	"testing"

	"github.com/Bose/minisentinel"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSentinel(t *testing.T) {
	m := miniredis.RunT(t)
	defer m.Close()

	masterName := "redis-master"

	replicaOpts := []minisentinel.Option{minisentinel.WithMasterName(masterName)}
	for range 2 {
		s := miniredis.RunT(t)
		defer s.Close()

		replicaOpts = append(replicaOpts, minisentinel.WithReplica(s))
	}

	s := minisentinel.NewSentinel(m, replicaOpts...)
	defer s.Close()
	require.NoError(t, s.Start())

	c, err := getSentinel(&Redis{
		MasterName: masterName,
		Addr:       s.Addr(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)
}
