package dartgen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	experimentalParser "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func TestGenDataEmbeddedStructs(t *testing.T) {
	tests := []struct {
		name        string
		types       string
		fields      []string
		properties  []string
		notContains []string
	}{
		{
			name:        "embedded struct",
			types:       "Common { ID int `json:\"id\"` }\nHealthRsp { Common }",
			fields:      []string{"final num id;"},
			properties:  []string{"id"},
			notContains: []string{"Common common", "Common.fromJson"},
		},
		{
			name:       "embedded struct declared after response",
			types:      "HealthRsp { Common }\nCommon { ID int `json:\"id\"` }",
			fields:     []string{"final num id;"},
			properties: []string{"id"},
		},
		{
			name:       "multiple embedded structs and own field",
			types:      "Common { ID int `json:\"id\"` }\nMetadata { Label string `json:\"label\"` }\nHealthRsp {\nCommon\nMetadata\nReady bool `json:\"ready\"`\n}",
			fields:     []string{"final num id;", "final String label;", "final bool ready;"},
			properties: []string{"id", "label", "ready"},
		},
		{
			name:        "nested embedded structs",
			types:       "Common { ID int `json:\"id\"` }\nMiddle {\nCommon\nLabel string `json:\"label\"`\n}\nHealthRsp { Middle }",
			fields:      []string{"final num id;", "final String label;"},
			properties:  []string{"id", "label"},
			notContains: []string{"Middle middle", "Middle.fromJson"},
		},
		{
			name:       "embedded pointer",
			types:      "Common { ID int `json:\"id\"` }\nHealthRsp { *Common }",
			fields:     []string{"final num id;"},
			properties: []string{"id"},
		},
		{
			name:        "own field shadows embedded field",
			types:       "Common { ID int `json:\"id\"` }\nHealthRsp {\nCommon\nID string `json:\"id\"`\n}",
			fields:      []string{"final String id;"},
			properties:  []string{"id"},
			notContains: []string{"final num id;"},
		},
		{
			name:        "shallower embedded field wins",
			types:       "Base { ID int `json:\"id\"` }\nCommon { Base }\nMetadata { ID string `json:\"id\"` }\nHealthRsp {\nCommon\nMetadata\n}",
			fields:      []string{"final String id;"},
			properties:  []string{"id"},
			notContains: []string{"final num id;"},
		},
		{
			name:        "named struct field stays nested",
			types:       "Common { ID int `json:\"id\"` }\nHealthRsp { Details Common `json:\"details\"` }",
			fields:      []string{"final Common details;"},
			properties:  []string{"details"},
			notContains: []string{"final num id;", "'id':"},
		},
		{
			name:        "empty embedded struct",
			types:       "Common {}\nHealthRsp { Common }",
			notContains: []string{"Common common", "Common.fromJson"},
		},
	}

	for _, legacy := range []bool{false, true} {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("legacy=%t/%s", legacy, tt.name), func(t *testing.T) {
				content := "syntax = \"v1\"\ntype (\n" + tt.types + "\nHealthReq {}\n)\n" +
					"service main {\n@handler GetHealth\nget /health (HealthReq) returns (HealthRsp)\n}\n"
				api, err := parser.ParseContent(content)
				require.NoError(t, err)
				dir := t.TempDir() + string(os.PathSeparator)
				require.NoError(t, genData(dir, api, legacy))
				data, err := os.ReadFile(filepath.Join(dir, "main.dart"))
				require.NoError(t, err)
				class := regexp.MustCompile(`(?s)class HealthRsp\s*\{(.*?)\n\}`).FindString(string(data))
				require.NotEmpty(t, class)
				class = strings.Join(strings.Fields(class), " ")
				for _, field := range tt.fields {
					require.Contains(t, class, field)
				}
				for _, property := range tt.properties {
					require.Contains(t, class, "this."+property)
					require.Contains(t, class, "m['"+property+"']")
					require.Contains(t, class, "'"+property+"':")
				}
				for _, unwanted := range tt.notContains {
					require.NotContains(t, class, unwanted)
				}
				if tt.name == "empty embedded struct" && !legacy {
					require.Contains(t, class, "HealthRsp();")
				}
			})
		}
	}
}

func TestGenApiEmbeddedPathParameters(t *testing.T) {
	for _, embedding := range []string{"", "Common", "*Common"} {
		t.Run("embedding="+embedding, func(t *testing.T) {
			member := "ID string `path:\"id\"`"
			if embedding != "" {
				member = embedding
			}
			content := "syntax = \"v1\"\ntype (\nCommon { ID string `path:\"id\"` }\n" +
				"HealthReq { " + member + " }\nHealthRsp {}\n)\n" +
				"service main {\n@handler GetHealth\nget /health/:id (HealthReq) returns (HealthRsp)\n}\n"
			api, err := parser.ParseContent(content)
			require.NoError(t, err)
			dir := t.TempDir() + string(os.PathSeparator)
			require.NoError(t, genData(dir, api, false))
			require.NoError(t, genApi(dir, api, false))
			data, err := os.ReadFile(filepath.Join(dir, "main.dart"))
			require.NoError(t, err)
			generated := strings.Join(strings.Fields(string(data)), " ")
			require.Contains(t, generated, "Future getHealth( String id,")
			require.Contains(t, generated, `apiGet("/health/${id}",`)
		})
	}
}

func TestGenDataRecursiveEmbeddedStruct(t *testing.T) {
	api, err := parser.ParseContent("syntax = \"v1\"\ntype HealthRsp {\nID int `json:\"id\"`\n*HealthRsp\n}\n")
	require.NoError(t, err)
	api.Service.Name = "main"
	err = genData(t.TempDir()+string(os.PathSeparator), api, false)
	require.ErrorContains(t, err, "cyclic embedded type")
}

func TestGenDataEmbeddedJSONTags(t *testing.T) {
	content := "syntax = \"v1\"\ntype (\nBase {\nID string `json:\"identifier\"`\n}\n" +
		"Common {\nBase\n}\nHealthRsp {\nCommon\n}\nHealthReq {}\n)\n" +
		"service main {\n@handler GetHealth\nget /health (HealthReq) returns (HealthRsp)\n}\n"
	for _, experimental := range []bool{false, true} {
		for _, legacy := range []bool{false, true} {
			t.Run(fmt.Sprintf("experimental=%t/legacy=%t", experimental, legacy), func(t *testing.T) {
				var api *spec.ApiSpec
				var err error
				if experimental {
					api, err = experimentalParser.Parse(filepath.Join(t.TempDir(), "test.api"), content)
				} else {
					api, err = parser.ParseContent(content)
				}
				require.NoError(t, err)
				dir := t.TempDir() + string(os.PathSeparator)
				require.NoError(t, genData(dir, api, legacy))
				data, err := os.ReadFile(filepath.Join(dir, "main.dart"))
				require.NoError(t, err)
				class := regexp.MustCompile(`(?s)class HealthRsp\s*\{(.*?)\n\}`).FindString(string(data))
				require.NotEmpty(t, class)
				class = strings.Join(strings.Fields(class), " ")
				require.Contains(t, class, "final String id;")
				require.Contains(t, class, "this.id")
				require.Contains(t, class, "m['identifier']")
				require.Contains(t, class, "'identifier': id")
				require.NotContains(t, class, "Common.fromJson")
			})
		}
	}
}

func TestGenDataAmbiguousEmbeddedFields(t *testing.T) {
	api, err := parser.ParseContent("syntax = \"v1\"\ntype (\n" +
		"Common { ID int `json:\"id\"` }\nMetadata { ID string `json:\"id\"` }\n" +
		"HealthRsp {\nCommon\nMetadata\n}\n)\n")
	require.NoError(t, err)
	api.Service.Name = "main"
	err = genData(t.TempDir()+string(os.PathSeparator), api, false)
	require.ErrorContains(t, err, "ambiguous embedded field ID")
}

func TestGenDataConflictingEmbeddedProperties(t *testing.T) {
	for _, field := range []string{"ID", "Id"} {
		t.Run(field, func(t *testing.T) {
			api, err := parser.ParseContent("syntax = \"v1\"\ntype (\n" +
				"Common { ID int `json:\"id\"` }\nHealthRsp {\nCommon\n" +
				field + " string `json:\"identifier\"`\n}\n)\n")
			require.NoError(t, err)
			api.Service.Name = "main"
			err = genData(t.TempDir()+string(os.PathSeparator), api, false)
			require.ErrorContains(t, err, "conflicting property names")
		})
	}
}
