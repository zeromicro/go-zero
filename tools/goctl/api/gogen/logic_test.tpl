package {{.pkgName}}

import (
    "context"
    "testing"

    {{.imports}}
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func Test{{.logic}}_{{.function}}(t *testing.T) {
    c := config.Config{}
    mockSvcCtx := svc.NewServiceContext(c)
    // init mock service context here

    tests := []struct {
        name       string
        ctx        context.Context
        setupMocks func()
        {{if .hasRequest}}req        *{{.requestType}}{{end}}
        wantErr    bool
        checkResp  func(*{{.responseType}}, error)
    }{
        {
            name: "successful",
            ctx:  context.Background(),
            setupMocks: func() {
                // Mock data for this test case
            },
            {{if .hasRequest}}req:  &{{.requestType}}{
                // TODO: init your request here
            },{{end}}
            wantErr: false,
            checkResp: func(resp *{{.responseType}}, err error) {
                // TODO: Add your check logic here
            },
        },
        {
            name: "response error",
            ctx:  context.Background(),
            setupMocks: func() {
                // mock data for this test case
            },
            {{if .hasRequest}}req:  &{{.requestType}}{
                // TODO: init your request here
            },{{end}}
            wantErr: true,
            checkResp: func(resp *{{.responseType}}, err error) {
                // TODO: Add your check logic here
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMocks()
            l := New{{.logic}}(tt.ctx, mockSvcCtx)
            resp, err := l.{{.function}}({{if .hasRequest}}tt.req{{end}})
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, resp)
            } else {
                require.NoError(t, err)
                assert.NotNil(t, resp)
                tt.checkResp(resp, err)
            }
        })
    }
}