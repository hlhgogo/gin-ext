package tracing

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestSpan_SpanContext(t *testing.T) {

	time.Sleep(time.Second * 3) // wait for jaeger initial

	span := Span{}

	opSpan := span.StartChildSpan("xxx")

	span = SpanFromOpentracing(opSpan)
	t.Logf("traceid: %s, span: %+v", span, opSpan)

}

func TestSpan_SpanText(t *testing.T) {

	var span Span = nil

	t.Logf(" span: %+v", !span.Empty())

}

func TestSpan_GoRoutinePanic(t *testing.T) {
	ctx := NewContext(context.TODO(), "aaa")
	for i := 0; i < 1001; i++ {
		go func() {
			NewContext(ctx, "aaa")
		}()
	}
	time.Sleep(5 * time.Second)
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), "aa")
	span := SpanFromContext(ctx)
	t.Log("span", span.AuthAccountID())
}

func TestKubenetesNamae(t *testing.T) {
	podNames := []string{
		"kube-apiserver-master01",
		"kube-controller-manager-master01",
		"kube-flannel-ds-4r9xf",
		"filebeat-5ls66",
		"alibaba-log-controller-b9d7bb4cb-cfhgc",
		"coredns-58cc8c89f4-cnst7",
		"",
	}

	for _, x := range podNames {
		names := strings.Split(x, "-")
		appName := strings.Join(names[:len(names)-1], "-")
		t.Log(x, appName)
	}
}
