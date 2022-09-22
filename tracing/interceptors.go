package tracing

import (
	"context"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Extract extracts the inbound HTTP request to obtain the parent span's context to ensure
// correct propagation of span context throughout the trace.
func Extract(r *http.Request, extHeaders ...string) (Span, error) {
	if r == nil {
		return nil, ErrNullRequest
	}
	if span := SpanFromContext(r.Context()); !span.Empty() {
		return span, nil
	}
	return extract(r.Header, extHeaders...), nil
}

func ExtractFromGrpc(ctx context.Context, extHeaders ...string) (Span, error) {
	headersIn, _ := metadata.FromIncomingContext(ctx)
	return extract(headersIn, extHeaders...), nil
}

// Inject injects the outbound HTTP request with the given span's context to ensure
// correct propagation of span context throughout the trace.
func Inject(r *http.Request, span Span) error {
	if r == nil {
		return ErrNullRequest
	}
	if span == nil {
		return nil
	}
	for k, v := range span {
		r.Header.Set(k, v)
	}
	return nil
}

func InjectToGrpc(ctx context.Context, span Span) context.Context {
	kv := []string{}
	for k, v := range span {
		kv = append(kv, k, v)
	}
	return metadata.AppendToOutgoingContext(ctx, kv...)
}

// UnaryClientInterceptor for passing incoming metadata to outgoing metadata
func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	// Take the incoming metadata and transfer it to the outgoing metadata
	if span := SpanFromContext(ctx); !span.Empty() {
		ctx = span.InjectToGrpc(ctx)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// StreamClientInterceptor returns a new streaming client interceptor for OpenTracing.
func StreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// Take the incoming metadata and transfer it to the outgoing metadata
	if span := SpanFromContext(ctx); !span.Empty() {
		ctx = span.InjectToGrpc(ctx)
	}
	return streamer(ctx, desc, cc, method, opts...)
}

// ContextWithSpan returns a new `context.Context` that holds a reference to
// the span. If span is nil, a new context without an active span is returned.
func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, span)
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `empty map` if no such `Span` could be found.
func SpanFromContext(ctx context.Context) Span {
	val := ctx.Value(activeSpanKey)
	if sp, ok := val.(Span); ok {
		// Fix panic when goroutine use same span
		newSpan := Span{}
		for k, v := range sp {
			newSpan.Set(k, v)
		}
		return newSpan
	}
	if span := opentracing.SpanFromContext(ctx); span != nil {
		return SpanFromOpentracing(span)
	}
	if span, err := ExtractFromGrpc(ctx); err == nil && span != nil {
		return span
	}
	return Span{}
}

func extract(headers map[string][]string, extHeaders ...string) Span {
	span := make(Span)
	for rawKey, rawValue := range headers {
		var value string
		if len(rawValue) > 0 {
			value = rawValue[0]
		}
		key := strings.ToLower(rawKey)
		if matchHeader(headersToPropagate, key) {
			span[key] = value
		} else if matchHeader(extHeaders, key) {
			span[key] = value
		} else if matchHeaderPrefix(applicationHeaderPrefix, key) {
			span[key] = value
		}
	}
	return span
}

func matchHeader(headers []string, key string) bool {
	for _, k := range headers {
		if k == key {
			return true
		}
	}
	return false
}

func matchHeaderPrefix(prefixes []string, key string) bool {
	for _, k := range prefixes {
		if strings.HasPrefix(key, k) {
			return true
		}
	}
	return false
}
