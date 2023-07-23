package internal

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestRequestParser_buildJsonRequestParser(t *testing.T) {
	parser, err := buildJsonRequestParser(map[string]any{"a": make(chan int)}, nil)
	assert.NotNil(t, err)
	assert.Nil(t, parser)
}

type badBody struct{}

func (badBody) Read([]byte) (int, error) { return 0, errors.New("something bad") }
func (badBody) Close() error             { return nil }
