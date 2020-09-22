package zrpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/zrpc/internal/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProxy(t *testing.T) {
	tests := []struct {
		name    string
		amount  float32
		res     *mock.DepositResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			"invalid request with negative amount",
			-1.11,
			nil,
			codes.InvalidArgument,
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"valid request with non negative amount",
			0.00,
			&mock.DepositResponse{Ok: true},
			codes.OK,
			"",
		},
	}

	proxy := NewProxy("foo", WithDialOption(grpc.WithInsecure()),
		WithDialOption(grpc.WithContextDialer(dialer())))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := proxy.TakeConn(context.Background())
			assert.Nil(t, err)
			cli := mock.NewDepositServiceClient(conn)
			request := &mock.DepositRequest{Amount: tt.amount}
			response, err := cli.Deposit(context.Background(), request)
			if response != nil {
				assert.True(t, len(response.String()) > 0)
				if response.GetOk() != tt.res.GetOk() {
					t.Error("response: expected", tt.res.GetOk(), "received", response.GetOk())
				}
			}
			if err != nil {
				if e, ok := status.FromError(err); ok {
					if e.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", e.Code())
					}
					if e.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", e.Message())
					}
				}
			}
		})
	}
}
