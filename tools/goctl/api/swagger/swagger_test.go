package swagger

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func Test_pathVariable2SwaggerVariable(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "/api/:id", expected: "/api/{id}"},
		{input: "/api/:id/details", expected: "/api/{id}/details"},
		{input: "/:version/api/:id", expected: "/{version}/api/{id}"},
		{input: "/api/v1", expected: "/api/v1"},
		{input: "/api/:id/:action", expected: "/api/{id}/{action}"},
	}

	for _, tc := range testCases {
		result := pathVariable2SwaggerVariable(testingContext(t), tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestArrayDefinitionsBug(t *testing.T) {
	// Test case for the bug where array of structs with useDefinitions
	// generates incorrect swagger JSON structure

	// Context with useDefinitions enabled
	ctx := Context{
		UseDefinitions: true,
	}

	// Create a test struct containing an array of structs
	testStruct := spec.DefineStruct{
		RawName: "TestStruct",
		Members: []spec.Member{
			{
				Name: "ArrayField",
				Type: spec.ArrayType{
					Value: spec.DefineStruct{
						RawName: "ItemStruct",
						Members: []spec.Member{
							{
								Name: "ItemName",
								Type: spec.PrimitiveType{RawName: "string"},
								Tag:  `json:"itemName"`,
							},
						},
					},
				},
				Tag: `json:"arrayField"`,
			},
		},
	}

	// Get properties from the struct
	properties, _ := propertiesFromType(ctx, testStruct)

	// Check that we have the array field
	assert.Contains(t, properties, "arrayField")
	arrayField := properties["arrayField"]

	// Verify the array field has correct structure
	assert.Equal(t, "array", arrayField.Type[0])
	
	// Check that we have items
	assert.NotNil(t, arrayField.Items, "Array should have items defined")
	assert.NotNil(t, arrayField.Items.Schema, "Array items should have schema")

	// The FIX: $ref should be inside items, not at schema level
	hasRef := arrayField.Ref.String() != ""
	assert.False(t, hasRef, "Schema level should NOT have $ref")
	
	// The $ref should be in the items
	hasItemsRef := arrayField.Items.Schema.Ref.String() != ""
	assert.True(t, hasItemsRef, "Items should have $ref")
	assert.Equal(t, "#/definitions/ItemStruct", arrayField.Items.Schema.Ref.String())

	// Verify there are no other properties in the items when using $ref
	assert.Nil(t, arrayField.Items.Schema.Properties, "Items with $ref should not have properties")
	assert.Empty(t, arrayField.Items.Schema.Required, "Items with $ref should not have required")
	assert.Empty(t, arrayField.Items.Schema.Type, "Items with $ref should not have type")
}

func TestParseExtensionValue(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "30", expected: float64(30)},
		{input: "3.14", expected: 3.14},
		{input: `"100/min"`, expected: "100/min"},
		{input: "`business`", expected: "business"},
		{input: "[\"alice\",\"bob\"]", expected: []interface{}{"alice", "bob"}},
		{input: "{\"foo\":\"bar\"}", expected: map[string]interface{}{"foo": "bar"}},
		{input: "`[\"alice\",\"bob\"]`", expected: []interface{}{"alice", "bob"}},
		{input: "`{\"foo\":\"bar\"}`", expected: map[string]interface{}{"foo": "bar"}},
		{input: "\"[\\\"alice\\\",\\\"bob\\\"]\"", expected: []interface{}{"alice", "bob"}},
		{input: "\"{\\\"foo\\\":\\\"bar\\\"}\"", expected: map[string]interface{}{"foo": "bar"}},
		{input: "plain", expected: "plain"},
	}
	for _, tt := range tests {
		actual := parseExtensionValue(tt.input)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestSpec2SwaggerWithRouteExtensions(t *testing.T) {
	src := `syntax = "v1"

type LoginRequest {
	Name string
}

type LoginResponse {
	Token string
}

service user {
	@doc "User Login"
	@x-log-enabled true
	@x-log-type business
	@x-sensitive true
	@x-rate-limit "100/min"
	@x-timeout 30
	@x-owners ` + "`" + `["alice","bob"]` + "`" + `
	@handler Login
	post /user/login (LoginRequest) returns (LoginResponse)
}`

	spec, err := parser.Parse("", src)
	assert.NoError(t, err)

	swagger, err := spec2Swagger(spec)
	assert.NoError(t, err)

	pathItem := swagger.Paths.Paths["/user/login"]
	assert.NotNil(t, pathItem.Post)
	op := pathItem.Post

	assert.Equal(t, true, op.Extensions["x-log-enabled"])
	assert.Equal(t, "business", op.Extensions["x-log-type"])
	assert.Equal(t, true, op.Extensions["x-sensitive"])
	assert.Equal(t, "100/min", op.Extensions["x-rate-limit"])
	assert.Equal(t, float64(30), op.Extensions["x-timeout"])
	assert.Equal(t, []interface{}{"alice", "bob"}, op.Extensions["x-owners"])

	// Also ensure the generated JSON really contains the x- fields.
	data, err := json.Marshal(swagger)
	assert.NoError(t, err)
	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))
	paths := m["paths"].(map[string]interface{})
	post := paths["/user/login"].(map[string]interface{})["post"].(map[string]interface{})
	assert.Equal(t, true, post["x-log-enabled"])
	assert.Equal(t, "business", post["x-log-type"])
}

func TestArrayWithoutDefinitions(t *testing.T) {
	// Test that arrays work correctly when useDefinitions is false
	ctx := Context{
		UseDefinitions: false, // This is the default
	}

	// Create the same test struct
	testStruct := spec.DefineStruct{
		RawName: "TestStruct",
		Members: []spec.Member{
			{
				Name: "ArrayField",
				Type: spec.ArrayType{
					Value: spec.DefineStruct{
						RawName: "ItemStruct",
						Members: []spec.Member{
							{
								Name: "ItemName",
								Type: spec.PrimitiveType{RawName: "string"},
								Tag:  `json:"itemName"`,
							},
						},
					},
				},
				Tag: `json:"arrayField"`,
			},
		},
	}

	properties, _ := propertiesFromType(ctx, testStruct)

	assert.Contains(t, properties, "arrayField")
	arrayField := properties["arrayField"]

	// Should be array type
	assert.Equal(t, "array", arrayField.Type[0])

	// Should have items with full schema, no $ref
	assert.NotNil(t, arrayField.Items)
	assert.NotNil(t, arrayField.Items.Schema)

	// Should NOT have $ref at schema level
	assert.Empty(t, arrayField.Ref.String(), "Schema should not have $ref when useDefinitions is false")

	// Should NOT have $ref in items either
	assert.Empty(t, arrayField.Items.Schema.Ref.String(), "Items should not have $ref when useDefinitions is false")

	// Should have full schema properties in items
	assert.Equal(t, "object", arrayField.Items.Schema.Type[0])
	assert.Contains(t, arrayField.Items.Schema.Properties, "itemName")
	assert.Equal(t, []string{"itemName"}, arrayField.Items.Schema.Required)
}
