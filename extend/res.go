package extend

import (
	"context"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	goError "github.com/go-errors/errors"
	"github.com/hlhgogo/config"
	athCtx "github.com/hlhgogo/gin-ext/context"
	"github.com/hlhgogo/gin-ext/errors"
	"github.com/hlhgogo/gin-ext/log"
	"net/http"
	"strings"
	"time"
)

// Res api response结构
type Res struct {
	Success bool        `json:"success"`
	Code    int         `json:"code,omitempty"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Debug   interface{} `json:"debug,omitempty"`
	TraceID string      `json:"traceId,omitempty"`
}

// ResDebug 响应带debug信息
type ResDebug struct {
	*Res
	Debug map[string]interface{}
}

// SendSuccess 返回成功结果
func SendSuccess(ctx *gin.Context, resData interface{}) {
	var res *Res
	res = successRes()

	res.Data = resData
	res.TraceID = athCtx.GetTraceId(ctx.Request.Context())

	if responseByte, err := json.Marshal(res); err == nil {
		log.InfoWithTrace(ctx.Request.Context(), "Response:%s", string(responseByte))
	}

	ctx.JSON(http.StatusOK, res)
}

// SendData 返回结果
// 根据err 是否为nil判断返回成功或失败
func SendData(ctx *gin.Context, resData interface{}, pErr error) {
	var res *Res
	var httpStatus = http.StatusOK
	if pErr != nil {
		res = failedRes()
		captureException(ctx.Request.Context(), pErr)
		if e, ok := pErr.(*errors.Err); ok {
			httpStatus = http.StatusInternalServerError
			if code := e.Code(); code != 0 {
				res.Code = code
			}
			if msg := e.Message(); msg != "" {
				res.Msg = msg
			}

			var debug interface{}
			if stack, ok := ctx.Get("Stack"); ok {
				debug = stack
			} else {
				if pErr != nil {
					switch err := pErr.(type) {
					case *goError.Error:
						debug = strings.Split(err.ErrorStack(), "\n")
					}
				}
			}

			if config.Get().App.ShowTrace {
				res.Debug = debug
			}

		} else if e, ok := pErr.(*errors.BadRequestError); ok {
			httpStatus = http.StatusBadRequest
			if code := e.Code(); code != 0 {
				res.Code = code
			}
			if msg := e.Message(); msg != "" {
				res.Msg = msg
			}
		} else if e, ok := pErr.(*errors.UnauthorizedError); ok {
			httpStatus = http.StatusUnauthorized
			if code := e.Code(); code != 0 {
				res.Code = code
			}
			if msg := e.Message(); msg != "" {
				res.Msg = msg
			}
		} else if e, ok := pErr.(*errors.ErrNotFoundError); ok {
			httpStatus = http.StatusNotFound
			if code := e.Code(); code != 0 {
				res.Code = code
			}
			if msg := e.Message(); msg != "" {
				res.Msg = msg
			}
		}
	} else {
		res = successRes()
	}

	res.Data = resData
	res.TraceID = athCtx.GetTraceId(ctx.Request.Context())

	if config.Get().App.ShowTrace {
		if stack, ok := ctx.Get("Stack"); ok {
			res.Debug = stack
		}
	}

	// 记录响应日志
	if responseByte, err := json.Marshal(res); err == nil {
		log.InfoWithTrace(ctx.Request.Context(), "Response:%s", string(responseByte))
	}

	ctx.JSON(httpStatus, res)
}

// newRes 新建
func newRes(success bool, code int, data interface{}) *Res {

	var msg string
	if value, ok := errors.ErrText[code]; ok {
		msg = value
	}
	return &Res{
		Success: success,
		Data:    data,
		Code:    code,
		Msg:     msg, // 原始，会被err中的msg替换 ，err中没有msg,会显示未定义
	}
}

// defaultRes 默认
func defaultRes() *Res {
	return newRes(true, errors.Success, struct{}{})
}

// failedRes 失败
func failedRes() *Res {
	return newRes(false, errors.ErrInternalServerError, struct{}{})
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

// successRes 成功
func successRes() *Res {
	return defaultRes()
}
