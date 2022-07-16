package gateway

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	metadataHeaderPrefix = "Grpc-Metadata-"
	metadataPrefix       = "gateway-"
)

func buildHeaders(header http.Header) []string {
	var headers []string

	for k, v := range header {
		if !strings.HasPrefix(k, metadataHeaderPrefix) {
			continue
		}

		key := fmt.Sprintf("%s%s", metadataPrefix, strings.TrimPrefix(k, metadataHeaderPrefix))
		for _, vv := range v {
			headers = append(headers, key+":"+vv)
		}
	}

	return headers
}
