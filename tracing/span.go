package tracing

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

const (
	HeaderRequestID     = "x-request-id"
	HeaderAuthAccountID = "x-auth-accountid"
)

type contextKey struct{}

var (
	headersToPropagate = []string{
		// All applications should propagate x-request-id. This header is
		// included in access log statements and is used for consistent trace
		// sampling and log sampling decisions in Istio.
		"x-request-id",

		// b3 trace headers. Compatible with Zipkin, OpenCensusAgent, and
		// Stackdriver Istio configurations. Commented out since they are
		// propagated by the OpenTracing tracer above.
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		// real ip
		"x-forwarded-for",
		"x-real-ip",
	}

	applicationHeaderPrefix = []string{
		"x-auth-",
	}

	ErrNullRequest = errors.New("original request is nil")

	activeSpanKey = contextKey{}
)

// Span represents propagated span identity and state
type Span map[string]string

func (span Span) Set(key, value string) {
	span[key] = value
}

func (span Span) Get(key string) string {
	return span[key]
}

func (span Span) RequestID() string {
	id := span.istioRequestID()
	if len(id) > 0 {
		return id
	}
	return span.traceID()
}

func (span Span) SpanID() string {
	return span.Get("x-b3-spanid")
}

func (span Span) traceID() string {
	return span.Get("x-b3-traceid")
}

func (span Span) istioRequestID() string {
	return span.Get(HeaderRequestID)
}

func (span Span) AuthAccountID() string {
	return span.Get(HeaderAuthAccountID)
}

func (span Span) SetAuthAccountID(accountID string) {
	span.Set(HeaderAuthAccountID, accountID)
}

func (span Span) IsIstio() bool {
	return istioOpen
}

func (span Span) Empty() bool {
	return len(span) == 0
}

func (span Span) Inject(r *http.Request) error {
	return Inject(r, span)
}

func (span Span) InjectToGrpc(context context.Context) context.Context {
	return InjectToGrpc(context, span)
}

func (span Span) ContextWithSpan(context context.Context) context.Context {
	return ContextWithSpan(context, span)
}

func (span Span) SpanContext() (opentracing.SpanContext, error) {
	var (
		traceID  jaeger.TraceID
		spanID   uint64
		parentID uint64
		sampled  = false
		baggage  = make(map[string]string)
		err      error
	)

	for key, value := range span {
		if key == "x-b3-traceid" {
			traceID, err = jaeger.TraceIDFromString(value)
		} else if key == "x-b3-parentspanid" {
			parentID, err = strconv.ParseUint(value, 16, 64)
		} else if key == "x-b3-spanid" {
			spanID, err = strconv.ParseUint(value, 16, 64)
		} else if key == "x-b3-sampled" && (value == "1" || value == "true") {
			sampled = true
		} else {
			baggage[key] = value
		}
		if err != nil {
			return nil, err
		}
	}

	spanctx := jaeger.NewSpanContext(
		traceID,
		jaeger.SpanID(spanID),
		jaeger.SpanID(parentID),
		sampled,
		baggage)

	return spanctx, nil
}

// NewContext 创建与公司有关的context，用于灰度链路传递
func NewContext(ctx context.Context, accountID string) context.Context {
	span := SpanFromContext(ctx)
	span.SetAuthAccountID(accountID)
	return ContextWithSpan(ctx, span)
}

// NewAccountIrrelevantContext 创建与公司信息无法关联的context
func NewAccountIrrelevantContext(ctx context.Context) context.Context {
	return NewContext(ctx, "-")
}

// CopyContext 将原始 context 链路信息复制到新 context 中
func CopyContext(new, old context.Context) context.Context {
	if new == nil {
		new = context.TODO()
	}
	span := SpanFromContext(old)
	return span.ContextWithSpan(new)
}

func (span Span) StartChildSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	parent, err := span.SpanContext()
	if err == nil {
		opts = append(opts, opentracing.ChildOf(parent))
	}
	newspan := opentracing.StartSpan(operationName, opts...)
	accountID := span.AuthAccountID()
	if len(accountID) > 0 {
		newspan.SetTag(HeaderAuthAccountID, accountID)
	}
	requestId := span.istioRequestID()
	if len(requestId) > 0 {
		newspan.SetTag("guid:x-request-id", span.istioRequestID())
	}
	tagOpentracingSpan(newspan)
	return newspan
}

func StartChildSpan(operationName string, span opentracing.Span, opts ...opentracing.StartSpanOption) opentracing.Span {
	opts = append(opts, opentracing.ChildOf(span.Context()))
	newspan := opentracing.StartSpan(operationName, opts...)
	tagOpentracingSpan(newspan)
	return newspan
}

func FinishSpan(span opentracing.Span) {
	span.Finish()
}

func StartNewSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	newspan := opentracing.StartSpan(operationName, opts...)
	tagOpentracingSpan(newspan)
	return newspan
}

func NewSpan(operationName string) Span {
	span := StartNewSpan(operationName)
	return SpanFromOpentracing(span)
}

func tagOpentracingSpan(span opentracing.Span) {
	spanCtx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return
	}
	span.SetTag("x-b3-traceid", spanCtx.TraceID().String())
}

func SpanFromOpentracing(span opentracing.Span) Span {
	spanCtx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return Span{}
	}

	newSpan := Span{}
	newSpan.Set("x-b3-traceid", spanCtx.TraceID().String())
	newSpan.Set("x-b3-parentspanid", spanCtx.ParentID().String())
	newSpan.Set("x-b3-spanid", spanCtx.SpanID().String())
	if spanCtx.IsSampled() {
		newSpan.Set("x-b3-sampled", "1")
	}

	spanCtx.ForeachBaggageItem(func(k, v string) bool {
		newSpan.Set(k, v)
		return true
	})

	return newSpan
}
