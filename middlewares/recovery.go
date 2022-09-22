package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hlhgogo/gin-ext/errors"
	"github.com/hlhgogo/gin-ext/extend"
	"github.com/hlhgogo/gin-ext/log"
	"runtime"
	"runtime/debug"
	"strings"
)

// PageNotFound 404PageHandle
func PageNotFound(c *gin.Context) {
	extend.SendData(c, nil, errors.NewErrNotFoundError())
}

// Recovery http request exception recovery
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stk := Stack(3, err)
				c.Set("Stack", stk)
				log.InfoMapWithTrace(c, stk, "Program Panic")
				switch err := err.(type) {
				case error:
					extend.SendData(c, nil, err)
				default:
					resErr := errors.NewErr("Unknown error")
					extend.SendData(c, nil, resErr)
				}
			}
		}()
		c.Next()
	}
}

// Stack get stack
func Stack(skip int, r interface{}) map[string]interface{} {

	var errorMessage string
	var errorStack []string
	var line int
	var lastFile string

	bufStack := debug.Stack()
	if len(bufStack) == 0 {
		return nil
	}

	if r != nil {
		errorMessage = fmt.Sprintf("%v", r)
	}

	stack := bytes.Split(bufStack, []byte("\n"))
	for _, bt := range stack {
		info := strings.TrimSpace(string(bt))
		if info != "" {
			errorStack = append(errorStack, info)
		}
	}

	_, lastFile, line, _ = runtime.Caller(skip)

	return gin.H{
		"error_file":  lastFile,
		"error_line":  line,
		"error_msg":   errorMessage,
		"error_trace": errorStack,
	}
}
