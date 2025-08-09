package internal

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func TestNewRequestParserNoVar(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithVars(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req = pathvar.WithVars(req, map[string]string{"a": "b"})
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserNoVarWithBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(`{"a": "b"}`))
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithNegativeContentLength(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(`{"a": "b"}`))
	req.ContentLength = -1
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithVarsWithBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(`{"a": "b"}`))
	req = pathvar.WithVars(req, map[string]string{"c": "d"})
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithVarsWithWrongBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(`{"a": "b"`))
	req = pathvar.WithVars(req, map[string]string{"c": "d"})
	parser, err := NewRequestParser(req, nil)
	assert.NotNil(t, err)
	assert.Nil(t, parser)
}

func TestNewRequestParserWithForm(t *testing.T) {
	req := httptest.NewRequest("GET", "/val?a=b", nil)
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithNilBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/val?a=b", http.NoBody)
	req.Body = nil
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithBadBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/val?a=b", badBody{})
	req.Body = badBody{}
	parser, err := NewRequestParser(req, nil)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestNewRequestParserWithBadForm(t *testing.T) {
	req := httptest.NewRequest("GET", "/val?a%1=b", http.NoBody)
	parser, err := NewRequestParser(req, nil)
	assert.NotNil(t, err)
	assert.Nil(t, parser)
}

func TestRequestParser_buildJsonRequestParserFromMap(t *testing.T) {
	parser, err := buildJsonRequestParserFromMap(map[string]any{"a": make(chan int)}, nil)
	assert.NotNil(t, err)
	assert.Nil(t, parser)
}

// mockAnyResolver is a simple implementation of jsonpb.AnyResolver for testing
type mockAnyResolver struct{}

func (m *mockAnyResolver) Resolve(typeUrl string) (proto.Message, error) {
	return nil, nil
}

func TestNewRequestParserWithIgnoreUnknownFields(t *testing.T) {
	// Create a concrete resolver for testing
	resolver := &mockAnyResolver{}

	// Test case 1: No body, no vars - should work with both true and false
	req1 := httptest.NewRequest("GET", "/", http.NoBody)
	parser1, err1 := NewRequestParser(req1, resolver)
	assert.Nil(t, err1)
	assert.NotNil(t, parser1)

	req2 := httptest.NewRequest("GET", "/", http.NoBody)
	parser2, err2 := NewRequestParser(req2, resolver)
	assert.Nil(t, err2)
	assert.NotNil(t, parser2)

	// Test case 2: With JSON body - tests the body parsing path
	req3 := httptest.NewRequest("POST", "/", strings.NewReader(`{"field": "value"}`))
	parser3, err3 := NewRequestParser(req3, resolver)
	assert.Nil(t, err3)
	assert.NotNil(t, parser3)

	req4 := httptest.NewRequest("POST", "/", strings.NewReader(`{"field": "value"}`))
	parser4, err4 := NewRequestParser(req4, resolver)
	assert.Nil(t, err4)
	assert.NotNil(t, parser4)
}

func TestNewRequestParserWithVarsAndIgnoreUnknownFields(t *testing.T) {
	resolver := &mockAnyResolver{}

	// Test with path variables and ignoreUnknownFields = true
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req = pathvar.WithVars(req, map[string]string{"a": "b"})
	parser, err := NewRequestParser(req, resolver)
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	// Test with path variables and ignoreUnknownFields = false
	req2 := httptest.NewRequest("GET", "/", http.NoBody)
	req2 = pathvar.WithVars(req2, map[string]string{"c": "d"})
	parser2, err2 := NewRequestParser(req2, resolver)
	assert.Nil(t, err2)
	assert.NotNil(t, parser2)
}

func TestNewRequestParserWithBodyAndIgnoreUnknownFields(t *testing.T) {
	resolver := &mockAnyResolver{}

	// Test with body and ignoreUnknownFields = true
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"a": "b"}`))
	parser, err := NewRequestParser(req, resolver)
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	// Test with body and ignoreUnknownFields = false
	req2 := httptest.NewRequest("POST", "/", strings.NewReader(`{"c": "d"}`))
	parser2, err2 := NewRequestParser(req2, resolver)
	assert.Nil(t, err2)
	assert.NotNil(t, parser2)
}

func TestNewRequestParserWithVarsBodyAndIgnoreUnknownFields(t *testing.T) {
	resolver := &mockAnyResolver{}

	// Test with both path variables and body, ignoreUnknownFields = true
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"a": "b"}`))
	req = pathvar.WithVars(req, map[string]string{"c": "d"})
	parser, err := NewRequestParser(req, resolver)
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	// Test with both path variables and body, ignoreUnknownFields = false
	req2 := httptest.NewRequest("POST", "/", strings.NewReader(`{"e": "f"}`))
	req2 = pathvar.WithVars(req2, map[string]string{"g": "h"})
	parser2, err2 := NewRequestParser(req2, resolver)
	assert.Nil(t, err2)
	assert.NotNil(t, parser2)
}

func TestBuildJsonRequestParserFromMapWithIgnoreUnknownFields(t *testing.T) {
	resolver := &mockAnyResolver{}

	// Test buildJsonRequestParserFromMap with ignoreUnknownFields = true
	data := map[string]any{"key": "value"}
	parser, err := buildJsonRequestParserFromMap(data, resolver)
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	// Test buildJsonRequestParserFromMap with ignoreUnknownFields = false
	parser2, err2 := buildJsonRequestParserFromMap(data, resolver)
	assert.Nil(t, err2)
	assert.NotNil(t, parser2)
}

func TestBuildJsonRequestParserWithUnknownFields(t *testing.T) {
	resolver := &mockAnyResolver{}

	// Test buildJsonRequestParserWithUnknownFields
	data := strings.NewReader(`{"test": "value"}`)
	parser, err := buildJsonRequestParserFromReader(data, resolver)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

type badBody struct{}

func (badBody) Read([]byte) (int, error) { return 0, errors.New("something bad") }
func (badBody) Close() error             { return nil }
