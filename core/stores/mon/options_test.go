package mon

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestWithRegistryForTwoRegisterType(t *testing.T) {
	opts := mopt.Client()

	// mongoDecimalEncoder allow user convert decimal.Decimal to primitive.Decimal128.
	var mongoDecimalEncoder bsoncodec.ValueEncoderFunc = func(ect bsoncodec.EncodeContext, w bsonrw.ValueWriter, value reflect.Value) error {
		// Use reflect, determine if it can be converted to decimal.Decimal.
		dec, ok := value.Interface().(decimal.Decimal)
		if !ok {
			return fmt.Errorf("value %v to encode is not of type decimal.Decimal", value)
		}

		// Convert decimal.Decimal to primitive.Decimal128.
		primDec, err := primitive.ParseDecimal128(dec.String())
		if err != nil {
			return fmt.Errorf("error converting decimal.Decimal %v to primitive.Decimal128: %v", dec, err)
		}
		return w.WriteDecimal128(primDec)
	}

	// mongoDecimalEncoder allow user convert primitive.Decimal128 to decimal.Decimal.
	var mongoDecimalDecoder bsoncodec.ValueDecoderFunc = func(ect bsoncodec.DecodeContext, r bsonrw.ValueReader, value reflect.Value) error {
		primDec, err := r.ReadDecimal128()
		if err != nil {
			return fmt.Errorf("error reading primitive.Decimal128 from ValueReader: %v", err)
		}

		// Convert primitive.Decimal128 to decimal.Decimal.
		dec, err := decimal.NewFromString(primDec.String())
		if err != nil {
			return fmt.Errorf("error converting primitive.Decimal128 %v to decimal.Decimal: %v", primDec, err)
		}

		// set value as decimal.Decimal type
		value.Set(reflect.ValueOf(dec))
		return nil
	}

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

	registerType := []RegisterType{
		{
			ValueType: reflect.TypeOf(decimal.Decimal{}),
			Encoder:   mongoDecimalEncoder,
			Decoder:   mongoDecimalDecoder,
		},
		{
			ValueType: reflect.TypeOf(time.Time{}),
			Encoder:   mongoDateTimeEncoder,
			Decoder:   mongoDateTimeDecoder,
		},
	}
	WithRegistry(registerType...)(opts)

	for _, v := range registerType {
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
