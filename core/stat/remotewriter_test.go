package stat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRemoteWriter(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").Reply(200).BodyString("foo")
	writer := NewRemoteWriter("http://foo.com")
	err := writer.Write(&StatReport{
		Name: "bar",
	})
	assert.Nil(t, err)
}

func TestRemoteWriterFail(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").Reply(503).BodyString("foo")
	writer := NewRemoteWriter("http://foo.com")
	err := writer.Write(&StatReport{
		Name: "bar",
	})
	assert.NotNil(t, err)
}

func TestRemoteWriterError(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").ReplyError(errors.New("foo"))
	writer := NewRemoteWriter("http://foo.com")
	err := writer.Write(&StatReport{
		Name: "bar",
	})
	assert.NotNil(t, err)
}
