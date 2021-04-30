package trace

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/core/timex"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

const (
	initSpanId  = "0"
	clientFlag  = "client"
	serverFlag  = "server"
	spanSepRune = '.'
)

var spanSep = string([]byte{spanSepRune})

// A Span is a calling span that connects caller and callee.
type Span struct {
	ctx           spanContext
	serviceName   string
	operationName string
	startTime     time.Time
	flag          string
	children      int
}

func newServerSpan(carrier Carrier, serviceName, operationName string) tracespec.Trace {
	traceId := stringx.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(traceIdKey)
		}
		return ""
	}, stringx.RandId)
	spanId := stringx.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(spanIdKey)
		}
		return ""
	}, func() string {
		return initSpanId
	})

	return &Span{
		ctx: spanContext{
			traceId: traceId,
			spanId:  spanId,
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timex.Time(),
		flag:          serverFlag,
	}
}

// Finish finishes the calling span.
func (s *Span) Finish() {
}

// Follow follows the tracing service and operation names in context.
func (s *Span) Follow(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId:  s.followSpanId(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timex.Time(),
		flag:          s.flag,
	}
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

// Fork forks the tracing service and operation names in context.
func (s *Span) Fork(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId:  s.forkSpanId(),
		},
		serviceName:   serviceName,
		operationName: operationName,
		startTime:     timex.Time(),
		flag:          clientFlag,
	}
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

// SpanId returns the span id.
func (s *Span) SpanId() string {
	return s.ctx.SpanId()
}

// TraceId returns the trace id.
func (s *Span) TraceId() string {
	return s.ctx.TraceId()
}

// Visit visits the span using fn.
func (s *Span) Visit(fn func(key, val string) bool) {
	s.ctx.Visit(fn)
}

func (s *Span) forkSpanId() string {
	s.children++
	return fmt.Sprintf("%s.%d", s.ctx.spanId, s.children)
}

func (s *Span) followSpanId() string {
	fields := strings.FieldsFunc(s.ctx.spanId, func(r rune) bool {
		return r == spanSepRune
	})
	if len(fields) == 0 {
		return s.ctx.spanId
	}

	last := fields[len(fields)-1]
	val, err := strconv.Atoi(last)
	if err != nil {
		return s.ctx.spanId
	}

	last = strconv.Itoa(val + 1)
	fields[len(fields)-1] = last

	return strings.Join(fields, spanSep)
}

// StartClientSpan starts the client span with given context, service and operation names.
func StartClientSpan(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	if span, ok := ctx.Value(tracespec.TracingKey).(*Span); ok {
		return span.Fork(ctx, serviceName, operationName)
	}

	return ctx, emptyNoopSpan
}

// StartServerSpan starts the server span with given context, carrier, service and operation names.
func StartServerSpan(ctx context.Context, carrier Carrier, serviceName, operationName string) (
	context.Context, tracespec.Trace) {
	span := newServerSpan(carrier, serviceName, operationName)
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}
