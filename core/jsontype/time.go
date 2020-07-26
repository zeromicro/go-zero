package jsontype

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
)

type MilliTime struct {
	time.Time
}

func (mt MilliTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.Milli())
}

func (mt *MilliTime) UnmarshalJSON(data []byte) error {
	var milli int64
	if err := json.Unmarshal(data, &milli); err != nil {
		return err
	} else {
		mt.Time = time.Unix(0, milli*int64(time.Millisecond))
		return nil
	}
}

func (mt MilliTime) GetBSON() (interface{}, error) {
	return mt.Time, nil
}

func (mt *MilliTime) SetBSON(raw bson.Raw) error {
	return raw.Unmarshal(&mt.Time)
}

func (mt MilliTime) Milli() int64 {
	return mt.UnixNano() / int64(time.Millisecond)
}
