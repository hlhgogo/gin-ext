package log

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	athCtx "github.com/hlhgogo/gin-ext/context"
	"github.com/sirupsen/logrus"
	"time"
)

// TraceWithTrace trace增加traceId
func TraceWithTrace(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelDebug)
	log.WithFields(getTraceField(ctx)).Trace(msg)
}

// DebugWithTrace debug增加traceId
func DebugWithTrace(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelDebug)
	log.WithFields(getTraceField(ctx)).Debug(msg)
}

// InfoWithTrace Info增加traceId
func InfoWithTrace(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelInfo)
	log.WithFields(getTraceField(ctx)).Info(msg)
}

// InfoMapWithTrace info增加map信息到日志
func InfoMapWithTrace(ctx context.Context, infos map[string]interface{}, format string, args ...interface{}) {
	fields := getTraceField(ctx)
	for k, v := range infos {
		fields[k] = v
	}
	log.WithFields(fields).Infof(format, args...)
}

// WarnWithTrace warn增加traceId
func WarnWithTrace(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelWarning)
	log.WithFields(getTraceField(ctx)).Warn(msg)
}

// WarnMapWithTrace warn增加map信息到日志
func WarnMapWithTrace(ctx context.Context, infos map[string]interface{}, format string, args ...interface{}) {
	fields := getTraceField(ctx)
	for k, v := range infos {
		fields[k] = v
	}
	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelWarning)
	log.WithFields(fields).Warnf(msg)
}

// ErrorWithTrace Error增加traceId
func ErrorWithTrace(ctx context.Context, err error, format string, args ...interface{}) {
	fields := getTraceField(ctx)
	fields["msg"] = err.Error()
	switch err := err.(type) {
	case *errors.Error:
		fields["stack"] = err.ErrorStack()
	default:
		newErr := errors.Wrap(fmt.Sprintf(format, args...), 1)
		fields["stack"] = newErr.ErrorStack()
	}

	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelError)
	captureException(ctx, err)
	log.WithFields(fields).Error(msg)
}

// ErrorMapWithTrace error增加map信息到日志
func ErrorMapWithTrace(ctx context.Context, infos map[string]interface{}, err error, format string, args ...interface{}) {
	fields := getTraceField(ctx)
	fields["msg"] = err.Error()
	for k, v := range infos {
		fields[k] = v
	}
	switch err := err.(type) {
	case *errors.Error:
		fields["stack"] = err.ErrorStack()
	default:
		newErr := errors.Wrap(fmt.Sprintf(format, args...), 1)
		fields["stack"] = newErr.ErrorStack()
	}

	msg := fmt.Sprintf(format, args...)
	addBreadcrumb(ctx, fmt.Sprintf(format, args...), sentry.LevelError)
	captureException(ctx, err)
	log.WithFields(fields).Error(msg)
}

// setBreadcrumb 增加一条sentry面板记录
func addBreadcrumb(ctx context.Context, msg string, level sentry.Level) {
	cv := athCtx.GetCtxValue(ctx)
	if cv == nil {
		return
	}
	if hub := cv.GetSentryHub(); hub != nil {
		hub.Scope().AddBreadcrumb(&sentry.Breadcrumb{
			Category: "logger",
			Message:  msg,
			Level:    level,
		}, 50)
	}
}

// captureException 上报异常
func captureException(ctx context.Context, err error) {
	cv := athCtx.GetCtxValue(ctx)
	if cv == nil {
		return
	}
	if hub := cv.GetSentryHub(); hub != nil {
		defer sentry.Flush(2 * time.Second)
		hub.CaptureException(err)
	}
}

// getTraceField 获取loggerField
func getTraceField(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{
		"type": Type,
	}

	requestId := athCtx.GetTraceId(ctx)
	fields["traceId"] = requestId

	return fields
}
