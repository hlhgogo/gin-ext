package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/hlhgogo/config"
	"strings"
)

func filterAlertWrapperFrames(frames []sentry.Frame) []sentry.Frame {
	filteredFrames := make([]sentry.Frame, 0, len(frames))
	for _, frame := range frames {
		if strings.Contains(frame.AbsPath, "ctx_log.go") {
			continue
		}
		if strings.Contains(frame.AbsPath, "recovery.go") {
			continue
		}
		if strings.Contains(frame.AbsPath, "res.go") {
			continue
		}
		filteredFrames = append(filteredFrames, frame)
	}
	return filteredFrames
}

func filterAlertWrapper(event *sentry.Event) *sentry.Event {
	for _, ex := range event.Exception {
		if ex.Stacktrace == nil {
			continue
		}
		ex.Stacktrace.Frames = filterAlertWrapperFrames(ex.Stacktrace.Frames)
	}
	// This interface is used when we extract stacktrace from caught strings, eg. in panics
	for _, th := range event.Threads {
		if th.Stacktrace == nil {
			continue
		}
		th.Stacktrace.Frames = filterAlertWrapperFrames(th.Stacktrace.Frames)
	}
	return event
}

// Load 初始化sentry配置
func Load() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:         config.Get().Sentry.Dsn,
		Environment: config.Get().Sentry.Environment,
		Release:     config.Get().Sentry.Release,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			return filterAlertWrapper(event)
		},
	})
}
