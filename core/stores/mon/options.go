package mon

import (
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const defaultTimeout = time.Second * 3

var (
	slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
	logMon        = syncx.ForAtomicBool(true)
	logSlowMon    = syncx.ForAtomicBool(true)
)

type (
	options = mopt.ClientOptions

	// Option defines the method to customize a mongo model.
	Option func(opts *options)

	// todo
	RegisterType struct {
		ValueType reflect.Type
		Encoder   bsoncodec.ValueEncoder
		Decoder   bsoncodec.ValueDecoder
	}
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

func defaultTimeoutOption() Option {
	return func(opts *options) {
		opts.SetTimeout(defaultTimeout)
	}
}

// WithTimeout set the mon client operation timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.SetTimeout(timeout)
	}
}

// WithRegistry sets todo
func WithRegistry(registerType ...RegisterType) Option {
	return func(opts *options) {
		registry := bson.NewRegistry()
		for _, v := range registerType {
			registry.RegisterTypeEncoder(v.ValueType, v.Encoder)
			registry.RegisterTypeDecoder(v.ValueType, v.Decoder)
		}
		opts.SetRegistry(registry)
	}
}
