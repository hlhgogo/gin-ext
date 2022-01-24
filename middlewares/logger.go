package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hlhgogo/config"
	athCtx "github.com/hlhgogo/gin-ext/context"
	"time"
)

// LoggerWithFormatter 格式化gin的日志输出
func LoggerWithFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency = param.Latency - param.Latency%time.Second
		}

		requestId := athCtx.GetTraceId(param.Request.Context())

		return fmt.Sprintf("%s [INFO] [%14s] [%13s] %5s | %s %-7s %s | %5s | %s %3d %s | %3s | %5s %s\n",
			config.Get().App.Name,
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			requestId,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Request.Proto,
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.Path,
			param.ErrorMessage,
		)
	})
}
