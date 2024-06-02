package mock

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DepositServer is used for mocking.
type DepositServer struct{}

// Deposit handles the deposit requests.
func (*DepositServer) Deposit(_ context.Context, req *DepositRequest) (*DepositResponse, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	time.Sleep(time.Duration(req.GetAmount()) * time.Millisecond)
	return &DepositResponse{Ok: true}, nil
}
