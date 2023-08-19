package internal

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	resText = `{"code":0,"data":{"pong":"pong"},"msg":"success"}`
	res     = Response{Pong: "pong"}
)

func TestEventHandler(t *testing.T) {

	h := NewEventHandler(io.Discard, nil)
	h.OnResolveMethod(nil)
	h.OnSendHeaders(nil)
	h.OnReceiveHeaders(nil)
	h.OnReceiveTrailers(status.New(codes.OK, ""), nil)
	h.OnReceiveResponse(nil)
	h.OnReceiveResponse(&res)
	assert.Equal(t, codes.OK, h.Status.Code())
}

func TestEventHandlerResponseTransform(t *testing.T) {

	w := &dataWriter{}
	h := NewEventHandler(w, nil, func(handler *EventHandler) {
		handler.RespHandler = func(writer io.Writer, status *status.Status, message proto.Message) {

			res, err := json.Marshal(map[string]interface{}{
				"code": status.Code(),
				"msg":  status.Message(),
				"data": message,
			})
			if err == nil {
				if _, we := writer.Write(res); we != nil {
					t.Error(we)
				}
			}
		}
	})
	h.OnReceiveResponse(&res)
	h.OnReceiveTrailers(status.New(codes.OK, "success"), nil)
	assert.Equal(t, codes.OK, h.Status.Code())
	assert.Equal(t, []byte(resText), w.data)
}

type dataWriter struct {
	data []byte
}

func (w *dataWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

type Response struct {
	Pong string `json:"pong"`
}

func (x *Response) Reset() {}
func (x *Response) String() string {
	d, _ := json.Marshal(x)
	return string(d)
}
func (*Response) ProtoMessage() {}
