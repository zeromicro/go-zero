package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"strings"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpc/internal"
	"github.com/zeromicro/go-zero/rest/internal/header"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var interceptors = []internal.Interceptor{
	internal.LogInterceptor,
}

// Do sends an HTTP request with the given arguments and returns an HTTP response.
// data is automatically marshal into a *httpRequest, typically it's defined in an API file.
func Do(ctx context.Context, method, url string, data any) (*http.Response, error) {
	req, err := buildRequest(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	return DoRequest(req)
}

// DoRequest sends an HTTP request and returns an HTTP response.
func DoRequest(r *http.Request) (*http.Response, error) {
	return request(r, defaultClient{})
}

type (
	client interface {
		do(r *http.Request) (*http.Response, error)
	}

	defaultClient struct{}
)

func (c defaultClient) do(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}

func buildFormQuery(u *nurl.URL, val map[string]any) string {
	query := u.Query()
	for k, v := range val {
		query.Add(k, fmt.Sprint(v))
	}

	return query.Encode()
}

func buildRequest(ctx context.Context, method, url string, data any) (*http.Request, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	var val map[string]map[string]any
	if data != nil {
		val, err = mapping.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	if err := fillPath(u, val[pathKey]); err != nil {
		return nil, err
	}

	var reader io.Reader
	jsonVars, hasJsonBody := val[jsonKey]
	if hasJsonBody {
		if method == http.MethodGet {
			return nil, ErrGetWithBody
		}

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(jsonVars); err != nil {
			return nil, err
		}

		reader = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = buildFormQuery(u, val[formKey])
	fillHeader(req, val[headerKey])
	if hasJsonBody {
		req.Header.Set(header.ContentType, header.JsonContentType)
	}

	return req, nil
}

func fillHeader(r *http.Request, val map[string]any) {
	for k, v := range val {
		r.Header.Add(k, fmt.Sprint(v))
	}
}

func fillPath(u *nurl.URL, val map[string]any) error {
	used := make(map[string]lang.PlaceholderType)
	fields := strings.Split(u.Path, slash)

	for i := range fields {
		field := fields[i]
		if len(field) > 0 && field[0] == colon {
			name := field[1:]
			ival, ok := val[name]
			if !ok {
				return fmt.Errorf("missing path variable %q", name)
			}
			value := fmt.Sprint(ival)
			if len(value) == 0 {
				return fmt.Errorf("empty path variable %q", name)
			}
			fields[i] = value
			used[name] = lang.Placeholder
		}
	}

	if len(val) != len(used) {
		for key := range used {
			delete(val, key)
		}

		unused := make([]string, 0, len(val))
		for key := range val {
			unused = append(unused, key)
		}

		return fmt.Errorf("more path variables are provided: %q", strings.Join(unused, ", "))
	}

	u.Path = strings.Join(fields, slash)
	return nil
}

func request(r *http.Request, cli client) (*http.Response, error) {
	ctx := r.Context()
	tracer := trace.TracerFromContext(ctx)
	propagator := otel.GetTextMapPropagator()

	spanName := r.URL.Path
	ctx, span := tracer.Start(
		ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(r)...),
	)
	defer span.End()

	respHandlers := make([]internal.ResponseHandler, len(interceptors))
	for i, interceptor := range interceptors {
		var h internal.ResponseHandler
		r, h = interceptor(r)
		respHandlers[i] = h
	}

	r = r.WithContext(ctx)
	propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

	resp, err := cli.do(r)
	for i := len(respHandlers) - 1; i >= 0; i-- {
		respHandlers[i](resp, err)
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return resp, err
	}

	span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(resp.StatusCode)...)
	span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(resp.StatusCode, oteltrace.SpanKindClient))

	return resp, err
}
