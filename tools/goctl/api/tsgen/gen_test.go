package tsgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

func TestGenWithInlineStructs(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	apiFile := filepath.Join(tmpDir, "test.api")

	// Write the test API file
	apiContent := `syntax = "v1"

info (
	title:   "Test ts generator"
	desc:    "Test inline struct handling"
	author:  "test"
	version: "v1"
)

// common pagination request
type PaginationReq {
	PageNum  int ` + "`form:\"pageNum\"`" + `
	PageSize int ` + "`form:\"pageSize\"`" + `
}

// base response
type BaseResp {
	Code int64  ` + "`json:\"code\"`" + `
	Msg  string ` + "`json:\"msg\"`" + `
}

// common req
type GetListCommonReq {
	Sth      string ` + "`form:\"sth\"`" + `
	PageNum  int    ` + "`form:\"pageNum\"`" + `
	PageSize int    ` + "`form:\"pageSize\"`" + `
}

// bad req to ts - inline struct with form tags
type GetListBadReq {
	Sth string ` + "`form:\"sth\"`" + `
	PaginationReq
}

// bad req to ts 2 - only inline struct with form tags
type GetListBad2Req {
	PaginationReq
}

// GetListResp - inline struct with json tags
type GetListResp {
	BaseResp
}

service test-api {
	@doc "common req"
	@handler getListCommon
	get /getListCommon (GetListCommonReq) returns (GetListResp)

	@doc "bad req"
	@handler getListBad
	get /getListBad (GetListBadReq) returns (GetListResp)

	@doc "bad req 2"
	@handler getListBad2
	get /getListBad2 (GetListBad2Req) returns (GetListResp)

	@doc "no req"
	@handler getListNoReq
	get /getListNoReq returns (GetListResp)
}`

	err := os.WriteFile(apiFile, []byte(apiContent), 0644)
	assert.NoError(t, err)

	// Parse the API file
	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	// Generate TypeScript files
	outputDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// Generate the files directly
	api.Service = api.Service.JoinPrefix()
	err = genRequest(outputDir)
	assert.NoError(t, err)
	err = genHandler(outputDir, ".", "webapi", api, false)
	assert.NoError(t, err)
	err = genComponents(outputDir, api)
	assert.NoError(t, err)

	// Read generated handler file
	handlerFile := filepath.Join(outputDir, "test.ts")
	handlerContent, err := os.ReadFile(handlerFile)
	assert.NoError(t, err)
	handler := string(handlerContent)

	// Read generated components file
	componentsFile := filepath.Join(outputDir, "testComponents.ts")
	componentsContent, err := os.ReadFile(componentsFile)
	assert.NoError(t, err)
	components := string(componentsContent)

	// Verify getListBad function signature and call
	assert.Contains(t, handler, "export function getListBad(params: components.GetListBadReqParams)")
	assert.Contains(t, handler, "return webapi.get<components.GetListResp>(`/getListBad`, params)")
	// Should NOT contain 4 arguments
	assert.NotContains(t, handler, "getListBad`, params, req, headers")

	// Verify getListBad2 function signature and call
	assert.Contains(t, handler, "export function getListBad2(params: components.GetListBad2ReqParams)")
	assert.Contains(t, handler, "return webapi.get<components.GetListResp>(`/getListBad2`, params)")
	// Should NOT reference non-existent headers
	assert.NotContains(t, handler, "GetListBad2ReqHeaders")

	// Verify getListCommon function signature and call
	assert.Contains(t, handler, "export function getListCommon(params: components.GetListCommonReqParams)")
	assert.Contains(t, handler, "return webapi.get<components.GetListResp>(`/getListCommon`, params)")

	// Verify getListNoReq function signature and call
	assert.Contains(t, handler, "export function getListNoReq()")
	assert.Contains(t, handler, "return webapi.get<components.GetListResp>(`/getListNoReq`)")

	// Verify GetListBadReqParams contains flattened fields
	assert.Contains(t, components, "export interface GetListBadReqParams")
	// Count occurrences of fields in GetListBadReqParams
	paramsStart := strings.Index(components, "export interface GetListBadReqParams")
	paramsEnd := strings.Index(components[paramsStart:], "}")
	paramsSection := components[paramsStart : paramsStart+paramsEnd]
	assert.Contains(t, paramsSection, "sth: string")
	assert.Contains(t, paramsSection, "pageNum: number")
	assert.Contains(t, paramsSection, "pageSize: number")

	// Verify GetListBad2ReqParams contains flattened fields from inline PaginationReq
	assert.Contains(t, components, "export interface GetListBad2ReqParams")
	params2Start := strings.Index(components, "export interface GetListBad2ReqParams")
	params2End := strings.Index(components[params2Start:], "}")
	params2Section := components[params2Start : params2Start+params2End]
	assert.Contains(t, params2Section, "pageNum: number")
	assert.Contains(t, params2Section, "pageSize: number")

	// Verify no empty Headers interfaces are generated
	assert.NotContains(t, components, "GetListBadReqHeaders")
	assert.NotContains(t, components, "GetListBad2ReqHeaders")

	// Verify GetListResp contains flattened fields from BaseResp
	assert.Contains(t, components, "export interface GetListResp")
	respStart := strings.Index(components, "export interface GetListResp")
	respEnd := strings.Index(components[respStart:], "}")
	respSection := components[respStart : respStart+respEnd]
	assert.Contains(t, respSection, "code: number")
	assert.Contains(t, respSection, "msg: string")
}
