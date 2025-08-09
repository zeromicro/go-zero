package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

// NewRequestParser creates a new request parser from the given http.Request and resolver.
func NewRequestParser(r *http.Request, resolver jsonpb.AnyResolver, ignoreUnknownFields bool) (grpcurl.RequestParser, error) {
	vars := pathvar.Vars(r)
	params, err := httpx.GetFormValues(r)
	if err != nil {
		return nil, err
	}

	for k, v := range vars {
		params[k] = v
	}

	body, ok := getBody(r)
	if !ok {
		return buildJsonRequestParserFromMap(params, resolver, ignoreUnknownFields)
	}

	if len(params) == 0 {
		if ignoreUnknownFields {
			return buildJsonRequestParserWithUnknownFields(body, resolver)
		}
		return buildJsonRequestParser(body, resolver)
	}

	m := make(map[string]any)
	if err := json.NewDecoder(body).Decode(&m); err != nil && err != io.EOF {
		return nil, err
	}

	for k, v := range params {
		m[k] = v
	}

	return buildJsonRequestParserFromMap(m, resolver, ignoreUnknownFields)
}

func buildJsonRequestParserFromMap(data map[string]any, resolver jsonpb.AnyResolver, ignoreUnknownFields bool) (grpcurl.RequestParser, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}

	if ignoreUnknownFields {
		return buildJsonRequestParserWithUnknownFields(&buf, resolver)
	}
	return buildJsonRequestParser(&buf, resolver)
}

// buildJsonRequestParser creates a JSON request parser with default settings
func buildJsonRequestParser(data io.Reader, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	return grpcurl.NewJSONRequestParser(data, resolver), nil
}

// buildJsonRequestParserWithUnknownFields creates a JSON request parser that ignores unknown fields
func buildJsonRequestParserWithUnknownFields(data io.Reader, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	unmarshaler := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
		AnyResolver:        resolver,
	}
	return grpcurl.NewJSONRequestParserWithUnmarshaler(data, unmarshaler), nil
}

func getBody(r *http.Request) (io.Reader, bool) {
	if r.Body == nil {
		return nil, false
	}

	if r.ContentLength == 0 {
		return nil, false
	}

	if r.ContentLength > 0 {
		return r.Body, true
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		return nil, false
	}

	if buf.Len() > 0 {
		return &buf, true
	}

	return nil, false
}
