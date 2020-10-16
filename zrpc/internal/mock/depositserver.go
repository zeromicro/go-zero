package mock

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DepositServer struct {
}

func (*DepositServer) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	return &DepositResponse{Ok: true}, nil
}
