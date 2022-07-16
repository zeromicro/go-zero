package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func buildJsonRequestParser(v interface{}, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}

	return grpcurl.NewJSONRequestParser(&buf, resolver), nil
}

func newRequestParser(r *http.Request, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	vars := pathvar.Vars(r)
	if len(vars) == 0 {
		return grpcurl.NewJSONRequestParser(r.Body, resolver), nil
	}

	if r.ContentLength == 0 {
		return buildJsonRequestParser(vars, resolver)
	}

	m := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}

	for k, v := range vars {
		m[k] = v
	}

	return buildJsonRequestParser(m, resolver)
}
