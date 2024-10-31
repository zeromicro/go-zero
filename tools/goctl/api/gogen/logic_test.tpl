package {{.pkgName}}

import (
    "context"
    "testing"

	{{.imports}}
	"github.com/stretchr/testify/assert"
)


func Test{{.logic}}_{{.function}}(t *testing.T) {
	c := config.Config{}
	mockSvcCtx := svc.NewServiceContext(c)

	tests := []struct {
		name       string
		ctx        context.Context
		setupMocks func()
		{{if .hasRequest}}req        *{{.requestType}}{{end}}
		wantErr    bool
		checkResp  func{{.responseType}}
	}{
		{
			name: "successful",
			ctx:  context.Background(),
			setupMocks: func() {
				// No setup needed for this test case
			},
			{{if .hasRequest}}req:  &{{.requestType}}{
                // init your request here
            },{{end}}
			wantErr: false,
			checkResp: func{{.responseType}} {
                // Add your check logic here
            },
		},
		{
			name: "response error",
			ctx:  context.Background(),
			setupMocks: func() {
				// No setup needed for this test case
			},
			{{if .hasRequest}}req:  &{{.requestType}}{
                // init your request here
            },{{end}}
			wantErr: true,
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
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				tt.checkResp(resp, err)
			}
		})
	}
}
