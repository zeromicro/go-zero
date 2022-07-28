package md

import (
	"context"
	"net/http"
	"strings"
)

// HeaderCarrier represents that the data in the header of http is converted into Metadata.
type HeaderCarrier http.Header

func (h HeaderCarrier) Extract(ctx context.Context) (context.Context, error) {
	metadata := FromContext(ctx)
	metadata = metadata.Clone()
	for k, v := range h {
		metadata.Append(strings.ToLower(k), v...)
	}

	return NewContext(ctx, metadata), nil
}

func (h HeaderCarrier) Inject(ctx context.Context) error {
	metadata := FromContext(ctx)
	for k, v := range metadata {
		h[strings.ToLower(k)] = v
	}

	return nil
}
