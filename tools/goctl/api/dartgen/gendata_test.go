package dartgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestGenDataClassListToJsonUsesToList(t *testing.T) {
	tests := []struct {
		name        string
		legacy      bool
		wantSnippet string
	}{
		{
			name:        "v2",
			wantSnippet: "items.map((i) => i?.toJson()).toList()",
		},
		{
			name:        "legacy",
			legacy:      true,
			wantSnippet: "'items': items.map((i) => i.toJson()).toList()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "data") + string(os.PathSeparator)

			if err := genData(dir, newClassListApiSpec(), tt.legacy); err != nil {
				t.Fatal(err)
			}

			data, err := os.ReadFile(filepath.Join(dir, "test.dart"))
			if err != nil {
				t.Fatal(err)
			}

			got := string(data)
			if !strings.Contains(got, tt.wantSnippet) {
				t.Fatalf("generated Dart does not contain %q:\n%s", tt.wantSnippet, got)
			}
		})
	}
}

func newClassListApiSpec() *spec.ApiSpec {
	return &spec.ApiSpec{
		Info: spec.Info{
			Title: "test",
		},
		Service: spec.Service{
			Name: "Test",
		},
		Types: []spec.Type{
			spec.DefineStruct{
				RawName: "GetTopUpProductsResponse",
				Members: []spec.Member{
					{
						Name: "Items",
						Type: spec.ArrayType{
							RawName: "[]*TopUpProductItem",
							Value: spec.PointerType{
								RawName: "*TopUpProductItem",
								Type: spec.DefineStruct{
									RawName: "TopUpProductItem",
								},
							},
						},
						Tag: "`json:\"items\"`",
					},
				},
			},
			spec.DefineStruct{
				RawName: "TopUpProductItem",
				Members: []spec.Member{
					{
						Name: "ProductId",
						Type: spec.PrimitiveType{
							RawName: "string",
						},
						Tag: "`json:\"productId\"`",
					},
				},
			},
		},
	}
}
