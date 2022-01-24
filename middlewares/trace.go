package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/hlhgogo/gin-ext/app"
	athCtx "github.com/hlhgogo/gin-ext/context"
	"github.com/hlhgogo/gin-ext/log"
	"github.com/satori/go.uuid"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {

		// bind request id
		requestId := c.Request.Header.Get("X-Request-Id")
		if requestId == "" {
			requestId = uuid.NewV4().String()
		}

		// save trace id to context
		commonValue := make(map[athCtx.CtxValueCommonKey]string)
		cv := athCtx.GetCtxValue(c.Request.Context())
		if cv != nil {
			commonValue = cv.GetCommonValue()
			commonValue[athCtx.CtxValueCommonKeyTraceID] = requestId
		}
		athValue := athCtx.GetCtxValue(c.Request.Context())
		athValue = athValue.SetCommonValue(commonValue)
		athContext, _ := athCtx.SetCtxValue(c.Request.Context(), athValue)
		c.Request = c.Request.WithContext(athContext)
		c.Header("Trace-Id", requestId)

		// ...
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// log with request info
		log.InfoMapWithTrace(c.Request.Context(), app.RequestInfo(c), "Request")

		c.Next()
	}
}
