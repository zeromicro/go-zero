package mock

import (
	"context"
	"errors"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func getRetryNum(ctx context.Context) int {
	md, b := metadata.FromIncomingContext(ctx)
	if b {
		n := md.Get("grpc-previous-rpc-attempts")

		if len(n) > 0 {
			if i, err := strconv.Atoi(n[0]); err == nil {
				return i
			}
		}
	}

	return 0
}

// DepositServer is used for mocking.
type DepositServer struct{}

// Deposit handles the deposit requests.
func (*DepositServer) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	retryNum := getRetryNum(ctx)

	// retry and error
	if req.GetAmount() == 2.11 {
		return nil, status.Errorf(codes.Unavailable, "need retry")
	}

	// retry and success
	if req.GetAmount() == 2.12 && retryNum < 1 {
		return nil, status.Errorf(codes.Unavailable, "need retry")
	}

	// not retry
	if req.GetAmount() == 2.13 {
		return nil, errors.New("not retry")
	}

	time.Sleep(time.Duration(req.GetAmount()) * time.Millisecond)
	return &DepositResponse{Ok: true}, nil
}
