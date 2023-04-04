package mock

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	retry = 0
)

func ResetRetry() {
	retry = 0
}

// DepositServer is used for mocking.
type DepositServer struct{}

// Deposit handles the deposit requests.
func (*DepositServer) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	// retry and error
	if req.GetAmount() == 2.11 {
		retry++
		return nil, status.Errorf(codes.Unavailable, "need retry")
	}

	// retry and success
	if req.GetAmount() == 2.12 && retry < 1 {
		retry++
		return nil, status.Errorf(codes.Unavailable, "need retry")
	}

	// not retry
	if req.GetAmount() == 2.13 {
		retry++
		return nil, errors.New("not retry")
	}

	time.Sleep(time.Duration(req.GetAmount()) * time.Millisecond)
	return &DepositResponse{Ok: true}, nil
}
