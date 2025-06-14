package sarama

import (
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func Test_requiredAcksFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    sarama.RequiredAcks
		wantErr bool
	}{
		{
			"none",
			args{s: "none"},
			sarama.NoResponse,
			false,
		},
		{
			"one",
			args{s: "one"},
			sarama.WaitForLocal,
			false,
		},
		{
			"all",
			args{s: "all"},
			sarama.WaitForAll,
			false,
		},
		{
			"invalid",
			args{s: "invalid"},
			sarama.NoResponse,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requiredAcksFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("requiredAcksFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("requiredAcksFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBalanceStrategy(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    sarama.BalanceStrategy
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"range",
			args{value: "range"},
			sarama.BalanceStrategyRange,
			assert.NoError,
		},
		{
			"roundrobin",
			args{value: "roundrobin"},
			sarama.BalanceStrategyRoundRobin,
			assert.NoError,
		},
		{
			"sticky",
			args{value: "sticky"},
			sarama.BalanceStrategySticky,
			assert.NoError,
		},
		{
			"invalid",
			args{value: "invalid"},
			nil,
			func(t assert.TestingT, err error, i ...any) bool {
				assert.EqualError(t, err, "kafka error: invalid BalanceStrategy: invalid")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBalanceStrategy(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("parseBalanceStrategy(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseBalanceStrategy(%v)", tt.args.value)
		})
	}
}

func Test_toSaramaConfig(t *testing.T) {
	type args struct {
		c types.ClientConfig
	}
	tests := []struct {
		name    string
		args    args
		want    func(*testing.T, *sarama.Config, ...any) bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty brokers",
			args: args{
				c: types.ClientConfig{
					Brokers: []string{},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.EqualError(t, err, "empty brokers")
				return true
			},
		},
		{
			name: "auto topic creation",
			args: args{
				c: types.ClientConfig{
					Brokers:                []string{"localhost:9092"},
					AllowAutoTopicCreation: true,
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.Equal(t, true, config.Metadata.AllowAutoTopicCreation)
				return true
			},
		},
		{
			name: "auth type password username empty",
			args: args{
				c: types.ClientConfig{
					Brokers:  []string{"localhost:9092"},
					AuthType: "password",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.EqualError(t, err, "username and password are required when using password auth type", i...)
				return true
			},
		},
		{
			name: "auth type password",
			args: args{
				c: types.ClientConfig{
					Brokers:      []string{"localhost:9092"},
					AuthType:     "password",
					SaslUsername: "username",
					SaslPassword: "password",
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.True(t, config.Net.SASL.Enable, i...)
				assert.Equal(t, "username", config.Net.SASL.User, i...)
				assert.Equal(t, "password", config.Net.SASL.Password, i...)
				assert.Equal(t, sarama.SASLMechanism(sarama.SASLTypePlaintext), config.Net.SASL.Mechanism, i...)
				return true
			},
		},
		{
			name: "tls enable",
			args: args{
				c: types.ClientConfig{
					Brokers:    []string{"localhost:9092"},
					TLSEnabled: true,
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.True(t, config.Net.TLS.Enable, i...)
				// add more tls config tests case
				assert.NotNil(t, config.Net.TLS.Config, i...)
				return true
			},
		},
		{
			name: "version",
			args: args{
				c: types.ClientConfig{
					Brokers: []string{"localhost:9092"},
					Version: "2.7.0",
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.Equal(t, sarama.V2_7_0_0, config.Version, i...)
				return true
			},
		},
		{
			name: "default version",
			args: args{
				c: types.ClientConfig{
					Brokers: []string{"localhost:9092"},
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.Equal(t, sarama.V2_4_0_0, config.Version, i...)
				return true
			},
		},
		{
			name: "modify clientId",
			args: args{
				c: types.ClientConfig{
					Brokers:  []string{"localhost:9092"},
					ClientId: "myapp",
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.Equal(t, "myapp", config.ClientID, i...)
				return true
			},
		},
		{
			name: "default clientId",
			args: args{
				c: types.ClientConfig{
					Brokers: []string{"localhost:9092"},
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.NotEmpty(t, config.ClientID, i...)
				return true
			},
		},
		{
			name: "req/resp size",
			args: args{
				c: types.ClientConfig{
					Brokers:         []string{"localhost:9092"},
					MaxRequestSize:  sarama.MaxRequestSize + 1,
					MaxResponseSize: sarama.MaxResponseSize + 1,
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.NotEmpty(t, config.ClientID, i...)
				return true
			},
		},
		{
			name: "saslMechanism SASLTypeSCRAMSHA256",
			args: args{
				c: types.ClientConfig{
					Brokers:       []string{"localhost:9092"},
					SaslMechanism: sarama.SASLTypeSCRAMSHA256,
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.NotEmpty(t, config.ClientID, i...)
				return true
			},
		},
		{
			name: "saslMechanism SASLTypeSCRAMSHA512",
			args: args{
				c: types.ClientConfig{
					Brokers:       []string{"localhost:9092"},
					SaslMechanism: sarama.SASLTypeSCRAMSHA512,
				},
			},
			want: func(t *testing.T, config *sarama.Config, i ...any) bool {
				assert.NotEmpty(t, config.ClientID, i...)
				return true
			},
		},
		{
			name: "invalid version",
			args: args{
				c: types.ClientConfig{
					Brokers: []string{"localhost:9092"},
					Version: "xxx",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				return true
			},
		},
		{
			name: "version is 2.3.0",
			args: args{
				c: types.ClientConfig{
					Brokers:   []string{"localhost:9092"},
					Version:   "2.3.0",
					AzEnabled: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toSaramaConfig(tt.args.c)
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("toSaramaConfig(%v)", tt.args.c)) {
					return
				}
			} else {
				assert.NoError(t, err, fmt.Sprintf("toSaramaConfig(%v)", tt.args.c))
			}
			if tt.want != nil {
				if !tt.want(t, got, fmt.Sprintf("toSaramaConfig(%v)", tt.args.c)) {
					return
				}
			}
		})
	}
}

func Test_compressionFromString(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    sarama.CompressionCodec
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"none",
			args{},
			sarama.CompressionNone,
			assert.NoError,
		},
		{
			"gzip",
			args{value: "gzip"},
			sarama.CompressionGZIP,
			assert.NoError,
		},
		{
			"snappy",
			args{value: "snappy"},
			sarama.CompressionSnappy,
			assert.NoError,
		},
		{
			"lz4",
			args{value: "lz4"},
			sarama.CompressionLZ4,
			assert.NoError,
		},
		{
			"invalid",
			args{value: "invalid"},
			0,
			func(t assert.TestingT, err error, i ...any) bool {
				assert.EqualError(t, err, "invalid compression: invalid")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compressionFromString(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("compressionFromString(%v)", tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "compressionFromString(%v)", tt.args.value)
		})
	}
}
