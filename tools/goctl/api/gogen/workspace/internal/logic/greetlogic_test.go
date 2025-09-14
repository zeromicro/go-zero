// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1-alpha

package logic

import (
	"context"
	"testing"

	"workspace/internal/config"
	"workspace/internal/svc"
	"workspace/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGreetLogic_Greet(t *testing.T) {
	c := config.Config{}
	mockSvcCtx := svc.NewServiceContext(c)
	// init mock service context here

	tests := []struct {
		name       string
		ctx        context.Context
		setupMocks func()
		req        *types.Request
		wantErr    bool
		checkResp  func(resp *types.Response, err error)
	}{
		{
			name: "response error",
			ctx:  context.Background(),
			setupMocks: func() {
				// mock data for this test case
			},
			req: &types.Request{
				// TODO: init your request here
			},
			wantErr: true,
			checkResp: func(resp *types.Response, err error) {
				// TODO: Add your check logic here
			},
		},
		{
			name: "successful",
			ctx:  context.Background(),
			setupMocks: func() {
				// Mock data for this test case
			},
			req: &types.Request{
				// TODO: init your request here
			},
			wantErr: false,
			checkResp: func(resp *types.Response, err error) {
				// TODO: Add your check logic here
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			l := NewGreetLogic(tt.ctx, mockSvcCtx)
			resp, err := l.Greet(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
			tt.checkResp(resp, err)
		})
	}
}
