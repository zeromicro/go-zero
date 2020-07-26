package handler

import (
	"time"

	"zero/stash/es"
	"zero/stash/filter"

	jsoniter "github.com/json-iterator/go"
)

const (
	timestampFormat = "2006-01-02T15:04:05.000Z"
	timestampKey    = "@timestamp"
)

type MessageHandler struct {
	writer  *es.Writer
	filters []filter.FilterFunc
}

func NewHandler(writer *es.Writer) *MessageHandler {
	return &MessageHandler{
		writer: writer,
	}
}

func (mh *MessageHandler) AddFilters(filters ...filter.FilterFunc) {
	for _, f := range filters {
		mh.filters = append(mh.filters, f)
	}
}

func (mh *MessageHandler) Consume(_, val string) error {
	m := make(map[string]interface{})
	if err := jsoniter.Unmarshal([]byte(val), &m); err != nil {
		return err
	}

	for _, proc := range mh.filters {
		if m = proc(m); m == nil {
			return nil
		}
	}

	bs, err := jsoniter.Marshal(m)
	if err != nil {
		return err
	}

	return mh.writer.Write(mh.getTime(m), string(bs))
}

func (mh *MessageHandler) getTime(m map[string]interface{}) time.Time {
	if ti, ok := m[timestampKey]; ok {
		if ts, ok := ti.(string); ok {
			if t, err := time.Parse(timestampFormat, ts); err == nil {
				return t
			}
		}
	}

	return time.Now()
}
