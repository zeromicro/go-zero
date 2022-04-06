package internal

import (
	"context"
	"net/http"
)

func Marshal(ctx context.Context, method, url string, data interface{}) (*http.Request, error) {
	panic("not implemented")
}
