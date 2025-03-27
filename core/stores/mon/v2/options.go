package mon

import (
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const defaultTimeout = time.Second * 3

var (
	slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
	logMon        = syncx.ForAtomicBool(true)
	logSlowMon    = syncx.ForAtomicBool(true)
)

type (
	// Option defines the method to customize a mongo model.
	Option func(opts *clientOptions)

	// TypeCodec is a struct that stores specific type Encoder/Decoder.
	TypeCodec struct {
		ValueType reflect.Type
		Encoder   bson.ValueEncoder
		Decoder   bson.ValueDecoder
	}

	clientOptions = options.ClientOptions
)

// DisableLog disables logging of mongo commands, includes info and slow logs.
func DisableLog() {
	logMon.Set(false)
	logSlowMon.Set(false)
}

// DisableInfoLog disables info logging of mongo commands, but keeps slow logs.
func DisableInfoLog() {
	logMon.Set(false)
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

// WithTimeout set the mon client operation timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *clientOptions) {
		opts.SetTimeout(timeout)
	}
}

// WithTypeCodec registers TypeCodecs to convert custom types.
func WithTypeCodec(typeCodecs ...TypeCodec) Option {
	return func(opts *clientOptions) {
		registry := bson.NewRegistry()
		for _, v := range typeCodecs {
			registry.RegisterTypeEncoder(v.ValueType, v.Encoder)
			registry.RegisterTypeDecoder(v.ValueType, v.Decoder)
		}
		opts.SetRegistry(registry)
	}
}

func defaultTimeoutOption() Option {
	return func(opts *clientOptions) {
		opts.SetTimeout(defaultTimeout)
	}
}
