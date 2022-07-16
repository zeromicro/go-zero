package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func newRequestParser(r *http.Request, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	vars := pathvar.Vars(r)
	if len(vars) == 0 {
		return grpcurl.NewJSONRequestParser(r.Body, resolver), nil
	}

	if r.ContentLength == 0 {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(vars); err != nil {
			return nil, err
		}

		return grpcurl.NewJSONRequestParser(&buf, resolver), nil
	}

	m := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}

	for k, v := range vars {
		m[k] = v
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}

	return grpcurl.NewJSONRequestParser(&buf, resolver), nil
}
