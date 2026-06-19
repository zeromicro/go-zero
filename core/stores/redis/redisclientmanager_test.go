package redis

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetClientTLSConfig(t *testing.T) {
	config := &tls.Config{
		ServerName:         "redis.example.com",
		InsecureSkipVerify: false,
	}
	r := newRedis("redis.example.com:6379", WithTLSConfig(config))

	client, err := getClient(r)
	require.NoError(t, err)
	actual := client.Options().TLSConfig
	require.NotNil(t, actual)
	assert.NotSame(t, config, actual)
	assert.Equal(t, "redis.example.com", actual.ServerName)
	assert.False(t, actual.InsecureSkipVerify)
}

func TestWithTLSUsesSecureDefaults(t *testing.T) {
	r := newRedis("redis.example.com:6379", WithTLS())

	require.NotNil(t, r.tlsConfig)
	assert.False(t, r.tlsConfig.InsecureSkipVerify)
}

func TestWithTLSConfigClonesConfig(t *testing.T) {
	config := &tls.Config{ServerName: "redis.example.com"}
	r := newRedis("redis.example.com:6379", WithTLSConfig(config))
	config.ServerName = "changed.example.com"

	require.NotNil(t, r.tlsConfig)
	assert.Equal(t, "redis.example.com", r.tlsConfig.ServerName)
}

func TestGetClientDoesNotShareAcrossTLSPolicies(t *testing.T) {
	const addr = "redis-tls-policy.example.com:6379"
	config := &tls.Config{ServerName: "redis.internal"}
	plainClient, err := getClient(newRedis(addr))
	require.NoError(t, err)
	tlsClient, err := getClient(newRedis(addr, WithTLS()))
	require.NoError(t, err)
	customTLSClient, err := getClient(newRedis(addr, WithTLSConfig(config)))
	require.NoError(t, err)
	reusedTLSClient, err := getClient(newRedis(addr, WithTLSConfig(config)))
	require.NoError(t, err)

	assert.NotSame(t, plainClient, tlsClient)
	assert.NotSame(t, tlsClient, customTLSClient)
	assert.Same(t, customTLSClient, reusedTLSClient)
	assert.Nil(t, plainClient.Options().TLSConfig)
	assert.NotNil(t, tlsClient.Options().TLSConfig)
	assert.Equal(t, "redis.internal", customTLSClient.Options().TLSConfig.ServerName)
}
