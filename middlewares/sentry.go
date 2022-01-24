package middlewares

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	athCtx "github.com/hlhgogo/gin-ext/context"
)

func Sentry() gin.HandlerFunc {
	return func(c *gin.Context) {
		hub := sentry.GetHubFromContext(c.Request.Context())
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
		}
		hub.Scope().SetRequest(c.Request)
		cv := athCtx.GetCtxValue(c.Request.Context())
		if cv != nil {
			commonValue := cv.GetCommonValue()
			if traceId, ok := commonValue[athCtx.CtxValueCommonKeyTraceID]; ok {
				hub.Scope().SetExtra("X-Request-Id", traceId)
			}
			athValue := cv.SetSentryHub(hub)
			athContext, _ := athCtx.SetCtxValue(c.Request.Context(), athValue)
			c.Request = c.Request.WithContext(athContext)
		}
		c.Next()
	}
}
