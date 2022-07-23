package md

import (
	"net/http"
	"strings"
)

var _ Carrier = (*HeaderCarrier)(nil)

// HeaderCarrier represents that the data in the header of http is converted into Metadata.
type HeaderCarrier http.Header

func (h HeaderCarrier) Carrier() (Metadata, error) {
	md := Metadata{}
	for k, v := range h {
		md[strings.ToLower(k)] = v
	}
	return md, nil
}
