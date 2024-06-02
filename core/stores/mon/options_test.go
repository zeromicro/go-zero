package mon

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func Test_defaultTimeoutOption(t *testing.T) {
	opts := mopt.Client()
	defaultTimeoutOption()(opts)
	assert.Equal(t, defaultTimeout, *opts.Timeout)
}

func TestWithTimeout(t *testing.T) {
	opts := mopt.Client()
	WithTimeout(time.Second)(opts)
	assert.Equal(t, time.Second, *opts.Timeout)
}

func TestDisableLog(t *testing.T) {
	assert.True(t, logMon.True())
	assert.True(t, logSlowMon.True())
	defer func() {
		logMon.Set(true)
		logSlowMon.Set(true)
	}()

	DisableLog()
	assert.False(t, logMon.True())
	assert.False(t, logSlowMon.True())
}

func TestDisableInfoLog(t *testing.T) {
	assert.True(t, logMon.True())
	assert.True(t, logSlowMon.True())
	defer func() {
		logMon.Set(true)
		logSlowMon.Set(true)
	}()

	DisableInfoLog()
	assert.False(t, logMon.True())
	assert.True(t, logSlowMon.True())
}

func TestWithRegistryForTimestampRegisterType(t *testing.T) {
	opts := mopt.Client()

	// mongoDateTimeEncoder allow user convert time.Time to primitive.DateTime.
	var mongoDateTimeEncoder bsoncodec.ValueEncoderFunc = func(ect bsoncodec.EncodeContext, w bsonrw.ValueWriter, value reflect.Value) error {
		// Use reflect, determine if it can be converted to time.Time.
		dec, ok := value.Interface().(time.Time)
		if !ok {
			return fmt.Errorf("value %v to encode is not of type time.Time", value)
		}
		return w.WriteDateTime(dec.Unix())
	}

	// mongoDateTimeEncoder allow user convert primitive.DateTime to time.Time.
	var mongoDateTimeDecoder bsoncodec.ValueDecoderFunc = func(ect bsoncodec.DecodeContext, r bsonrw.ValueReader, value reflect.Value) error {
		primTime, err := r.ReadDateTime()
		if err != nil {
			return fmt.Errorf("error reading primitive.DateTime from ValueReader: %v", err)
		}
		value.Set(reflect.ValueOf(time.Unix(primTime, 0)))
		return nil
	}

	codecs := []TypeCodec{
		{
			ValueType: reflect.TypeOf(time.Time{}),
			Encoder:   mongoDateTimeEncoder,
			Decoder:   mongoDateTimeDecoder,
		},
	}
	WithTypeCodec(codecs...)(opts)

	for _, v := range codecs {
		// Validate Encoder
		enc, err := opts.Registry.LookupEncoder(v.ValueType)
		if err != nil {
			t.Fatal(err)
		}
		if assert.ObjectsAreEqual(v.Encoder, enc) {
			t.Errorf("Encoder got from Registry: %v, but want: %v", enc, v.Encoder)
		}

		// Validate Decoder
		dec, err := opts.Registry.LookupDecoder(v.ValueType)
		if err != nil {
			t.Fatal(err)
		}
		if assert.ObjectsAreEqual(v.Decoder, dec) {
			t.Errorf("Decoder got from Registry: %v, but want: %v", dec, v.Decoder)
		}
	}
}
