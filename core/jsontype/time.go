package jsontype

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
)

// MilliTime represents time.Time that works better with mongodb.
type MilliTime struct {
	time.Time
}

// MarshalJSON marshals mt to json bytes.
func (mt MilliTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.Milli())
}

// UnmarshalJSON unmarshals data into mt.
func (mt *MilliTime) UnmarshalJSON(data []byte) error {
	var milli int64
	if err := json.Unmarshal(data, &milli); err != nil {
		return err
	}

	mt.Time = time.Unix(0, milli*int64(time.Millisecond))
	return nil
}

// GetBSON returns BSON base on mt.
func (mt MilliTime) GetBSON() (interface{}, error) {
	return mt.Time, nil
}

// SetBSON sets raw into mt.
func (mt *MilliTime) SetBSON(raw bson.Raw) error {
	return raw.Unmarshal(&mt.Time)
}

// Milli returns milliseconds for mt.
func (mt MilliTime) Milli() int64 {
	return mt.UnixNano() / int64(time.Millisecond)
}
