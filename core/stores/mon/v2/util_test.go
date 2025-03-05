package mon

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	"github.com/zeromicro/go-zero/core/timex"
)

func TestFormatAddrs(t *testing.T) {
	tests := []struct {
		addrs  []string
		expect string
	}{
		{
			addrs:  []string{"a", "b"},
			expect: "a,b",
		},
		{
			addrs:  []string{"a", "b", "c"},
			expect: "a,b,c",
		},
		{
			addrs:  []string{},
			expect: "",
		},
		{
			addrs:  nil,
			expect: "",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expect, FormatAddr(test.addrs))
	}
}

func Test_logDuration(t *testing.T) {
	buf := logtest.NewCollector(t)

	buf.Reset()
	logDuration(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil)
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slow")

	buf.Reset()
	logDuration(context.Background(), "foo", "bar", timex.Now(), nil)
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	logDuration(context.Background(), "foo", "bar", timex.Now(), errors.New("bar"))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "fail")

	defer func() {
		logMon.Set(true)
		logSlowMon.Set(true)
	}()

	buf.Reset()
	DisableInfoLog()
	logDuration(context.Background(), "foo", "bar", timex.Now(), nil)
	assert.Empty(t, buf.String())

	buf.Reset()
	logDuration(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil)
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slow")

	buf.Reset()
	DisableLog()
	logDuration(context.Background(), "foo", "bar", timex.Now(), nil)
	assert.Empty(t, buf.String())

	buf.Reset()
	logDuration(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil)
	assert.Empty(t, buf.String())

	buf.Reset()
	logDuration(context.Background(), "foo", "bar", timex.Now(), errors.New("bar"))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "fail")
}

func Test_logDurationWithDoc(t *testing.T) {
	buf := logtest.NewCollector(t)
	buf.Reset()

	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil, make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slow")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil, "{'json': ''}")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slow")
	assert.Contains(t, buf.String(), "json")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), nil, make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), nil, "{'json': ''}")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "json")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), errors.New("bar"), make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "fail")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), errors.New("bar"), "{'json': ''}")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "fail")
	assert.Contains(t, buf.String(), "json")

	defer func() {
		logMon.Set(true)
		logSlowMon.Set(true)
	}()

	buf.Reset()
	DisableInfoLog()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), nil, make(chan int))
	assert.Empty(t, buf.String())

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), nil, "{'json': ''}")
	assert.Empty(t, buf.String())

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil, make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slow")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil, "{'json': ''}")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slow")
	assert.Contains(t, buf.String(), "json")

	buf.Reset()
	DisableLog()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), nil, make(chan int))
	assert.Empty(t, buf.String())

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), nil, "{'json': ''}")
	assert.Empty(t, buf.String())

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil, make(chan int))
	assert.Empty(t, buf.String())

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now()-slowThreshold.Load()*2, nil, "{'json': ''}")
	assert.Empty(t, buf.String())

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), errors.New("bar"), make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "fail")

	buf.Reset()
	logDurationWithDocs(context.Background(), "foo", "bar", timex.Now(), errors.New("bar"), "{'json': ''}")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "fail")
	assert.Contains(t, buf.String(), "json")
}
