package trace

import (
	"net/http"

	"google.golang.org/grpc/metadata"
)

const (
	// HttpFormat means http carrier format.
	HttpFormat = iota
	// GrpcFormat means grpc carrier format.
	GrpcFormat
)

var (
	emptyHttpPropagator httpPropagator
	emptyGrpcPropagator grpcPropagator
)

type (
	// Propagator interface wraps the Extract and Inject methods.
	Propagator interface {
		Extract(carrier interface{}) (Carrier, error)
		Inject(carrier interface{}) (Carrier, error)
	}

	httpPropagator struct{}
	grpcPropagator struct{}
)

func (h httpPropagator) Extract(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(http.Header); ok {
		return httpCarrier(c), nil
	}

	return nil, ErrInvalidCarrier
}

func (h httpPropagator) Inject(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(http.Header); ok {
		return httpCarrier(c), nil
	}

	return nil, ErrInvalidCarrier
}

func (g grpcPropagator) Extract(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(metadata.MD); ok {
		return grpcCarrier(c), nil
	}

	return nil, ErrInvalidCarrier
}

func (g grpcPropagator) Inject(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(metadata.MD); ok {
		return grpcCarrier(c), nil
	}

	return nil, ErrInvalidCarrier
}

// Extract extracts tracing information from carrier with given format.
func Extract(format, carrier interface{}) (Carrier, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Extract(carrier)
		} else if v == GrpcFormat {
			return emptyGrpcPropagator.Extract(carrier)
		}
	}

	return nil, ErrInvalidCarrier
}

// Inject injects tracing information into carrier with given format.
func Inject(format, carrier interface{}) (Carrier, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Inject(carrier)
		} else if v == GrpcFormat {
			return emptyGrpcPropagator.Inject(carrier)
		}
	}

	return nil, ErrInvalidCarrier
}
