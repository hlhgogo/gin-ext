package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"reflect"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandString A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func InArray(val interface{}, array interface{}) (index int) {
	kind := reflect.TypeOf(array).Kind()
	values := reflect.ValueOf(array)

	if kind == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return i
			}
		}
	}

	return -1
}

func RequestInfo(c *gin.Context) map[string]interface{} {
	var (
		body string
		err  error
	)

	bodyByte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		body = ""
	} else {
		body = fmt.Sprintf("%v", string(bodyByte))
	}

	return map[string]interface{}{
		"headers": c.Request.Header,
		"query":   c.Request.URL.Query(),
		"body":    body,
	}
}

func Trace(c *gin.Context) map[string]interface{} {
	debugObj := gin.H{
		"request": RequestInfo(c),
	}
	if stack, exits := c.Get("Stack"); exits {
		debugObj["stack"] = stack
	}
	return debugObj
}
