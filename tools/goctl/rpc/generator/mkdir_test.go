package generator

import (
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

func TestServiceNameDetermination(t *testing.T) {
	tests := []struct {
		name             string
		protoName        string
		packageName      string
		hasPackage       bool
		nameFromFilename bool
		expectedName     string
	}{
		{
			name:             "default uses package name when available",
			protoName:        "user.proto",
			packageName:      "userservice",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "userservice",
		},
		{
			name:             "flag enabled uses filename instead of package",
			protoName:        "user.proto",
			packageName:      "userservice",
			hasPackage:       true,
			nameFromFilename: true,
			expectedName:     "user",
		},
		{
			name:             "fallback to filename when package is empty",
			protoName:        "order.proto",
			packageName:      "",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "order",
		},
		{
			name:             "fallback to filename when package is nil",
			protoName:        "product.proto",
			packageName:      "",
			hasPackage:       false,
			nameFromFilename: false,
			expectedName:     "product",
		},
		{
			name:             "flag enabled with nil package uses filename",
			protoName:        "catalog.proto",
			packageName:      "",
			hasPackage:       false,
			nameFromFilename: true,
			expectedName:     "catalog",
		},
		{
			name:             "handles proto file with complex name",
			protoName:        "user-service.proto",
			packageName:      "user",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "user",
		},
		{
			name:             "handles proto file with complex name and flag",
			protoName:        "user-service.proto",
			packageName:      "user",
			hasPackage:       true,
			nameFromFilename: true,
			expectedName:     "user-service",
		},
		{
			name:             "nil context uses package name",
			protoName:        "account.proto",
			packageName:      "accountpb",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "accountpb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the proto struct
			p := parser.Proto{
				Name: tt.protoName,
			}
			if tt.hasPackage {
				p.Package = parser.Package{
					Package: &proto.Package{
						Name: tt.packageName,
					},
				}
			}

			// Build the context
			var ctx *ZRpcContext
			if tt.name != "nil context uses package name" {
				ctx = &ZRpcContext{
					NameFromFilename: tt.nameFromFilename,
				}
			}

			// Call the helper function to determine service name
			serviceName := determineServiceName(p, ctx)
			assert.Equal(t, tt.expectedName, serviceName)
		})
	}
}

func TestServiceNameWithNilContext(t *testing.T) {
	p := parser.Proto{
		Name: "test.proto",
		Package: parser.Package{
			Package: &proto.Package{
				Name: "testpkg",
			},
		},
	}

	// nil context should use package name
	serviceName := determineServiceName(p, nil)
	assert.Equal(t, "testpkg", serviceName)
}

func TestServiceNameFallbackWithEmptyPackage(t *testing.T) {
	p := parser.Proto{
		Name: "myservice.proto",
		Package: parser.Package{
			Package: &proto.Package{
				Name: "", // empty package name
			},
		},
	}

	ctx := &ZRpcContext{
		NameFromFilename: false,
	}

	// Should fall back to filename when package name is empty
	serviceName := determineServiceName(p, ctx)
	assert.Equal(t, "myservice", serviceName)
}
