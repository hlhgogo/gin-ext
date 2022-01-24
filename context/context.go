package context

import (
	"context"
	"github.com/getsentry/sentry-go"
)

// CtxValueCommonKey ctx common value key
type CtxValueCommonKey string

// 常量定义
const (
	// CtxValueCommonKeyTraceID traceId
	CtxValueCommonKeyTraceID CtxValueCommonKey = "traceId"
)

// CtxValueKey ctx value key
type CtxValueKey string

// 常量定义
const (
	// CtxValueKeyV1 ...
	CtxValueKeyV1 CtxValueKey = "cfCtxValue"
)

// CtxValue 上下文内容
type CtxValue struct {
	common    map[CtxValueCommonKey]string
	sentryHub *sentry.Hub
}

// GetCommonValue get common value
func (v *CtxValue) GetCommonValue() map[CtxValueCommonKey]string {
	return v.common
}

// GetSentryHub 获取GetSentryHub
func (v *CtxValue) GetSentryHub() *sentry.Hub {
	return v.sentryHub
}

// SetCommonValue set common value
func (v *CtxValue) SetCommonValue(cv map[CtxValueCommonKey]string) *CtxValue {
	v.common = cv
	return v
}

// SetSentryHub set sentry hub
func (v *CtxValue) SetSentryHub(hub *sentry.Hub) *CtxValue {
	v.sentryHub = hub
	return v
}

// NewCtxValue 创建cxt value
func NewCtxValue(common map[CtxValueCommonKey]string) *CtxValue {
	return newCtxValue(common)
}

// newCtxValue 创建新的ctx value
func newCtxValue(common map[CtxValueCommonKey]string) *CtxValue {
	if common == nil {
		common = map[CtxValueCommonKey]string{}
	}
	return &CtxValue{
		common: common,
	}
}

// GetCtxValue 获取ctx value
func GetCtxValue(ctx context.Context) *CtxValue {
	cv := ctx.Value(CtxValueKeyV1)
	if v, ok := cv.(*CtxValue); ok {
		if v.common == nil {
			v.common = map[CtxValueCommonKey]string{}
		}
		return v
	}

	return &CtxValue{map[CtxValueCommonKey]string{}, nil}
}

// GetTraceId 获取traceId
func GetTraceId(ctx context.Context) string {
	cv := GetCtxValue(ctx)
	commonValue := cv.GetCommonValue()
	traceId := ""
	if value, ok := commonValue[CtxValueCommonKeyTraceID]; ok {
		traceId = value
	}

	return traceId
}

// SetCtxValue 设置ctx value
func SetCtxValue(ctx context.Context, value *CtxValue) (context.Context, *CtxValue) {
	ctx = context.WithValue(ctx, CtxValueKeyV1, value)
	return ctx, value
}
