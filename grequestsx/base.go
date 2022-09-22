package grequestsx

import (
	"github.com/hlhgogo/gin-ext/log"
	"github.com/hlhgogo/gin-ext/tracing"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
	"time"
)

type Flags struct {
	// EnableLog 记录请求和响应内容到日志
	EnableLog bool
	// DisableTrace 调用三方服务时，删除trace信息
	DisableTrace bool
}

// Get takes 2 parameters and returns a Response Struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Get(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("GET", url, ro, flags...)
}

// Put takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Put(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("PUT", url, ro, flags...)
}

// Patch takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Patch(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("PATCH", url, ro, flags...)
}

// Delete takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Delete(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("DELETE", url, ro, flags...)
}

// Post takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Post(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("POST", url, ro, flags...)
}

// Head takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Head(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("HEAD", url, ro, flags...)
}

// Options takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Options(url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest("OPTIONS", url, ro, flags...)
}

// Req takes 3 parameters and returns a Response Struct. These three options are:
//	1. A verb
// 	2. A URL
// 	3. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Req(verb string, url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	return DoRegularRequest(verb, url, ro, flags...)
}

// DoRegularRequest adds generic test functionality
func DoRegularRequest(requestVerb, url string, ro *grequests.RequestOptions, flags ...Flags) (*grequests.Response, error) {
	if ro == nil {
		ro = &grequests.RequestOptions{}
	}

	flag := Flags{}
	if len(flags) > 0 {
		flag = flags[0]
	}

	if !flag.DisableTrace {
		if ro.Context != nil {
			span := tracing.SpanFromContext(ro.Context)
			for k, v := range span {
				if ro.Headers == nil {
					ro.Headers = map[string]string{}
				}
				ro.Headers[k] = v
			}
			if span.Empty() {
				log.Warnf("[warning] mesher requested url has no tracing info: %s", url)
			}
		} else {
			log.Warnf("[warning] mesher requested url has no context info: %s", url)
		}
	}

	start := time.Now()
	response, err := grequests.DoRegularRequest(requestVerb, url, ro)
	elapsed := time.Since(start)

	if !flag.EnableLog {
		return response, err
	}

	// request logs
	reqLogFields := logrus.Fields{
		"url":    url,
		"method": requestVerb,
	}

	if len(ro.Headers) > 0 {
		reqLogFields["headers"] = ro.Headers
	}

	if len(ro.Data) > 0 {
		reqLogFields["data"] = ro.Data
	}
	if len(ro.Params) > 0 {
		reqLogFields["params"] = ro.Params
	}
	if ro.JSON != nil {
		reqLogFields["json"] = ro.JSON
	}
	if ro.XML != nil {
		reqLogFields["xml"] = ro.XML
	}

	if err != nil {
		log.ErrorFieldsWithTrace(ro.Context, logrus.Fields{
			"request": reqLogFields,
			"elapsed": elapsed.Seconds(),
		}, err, "requestsx request failed")
		return response, err
	}

	respLogFields := logrus.Fields{}
	respLogFields["status"] = response.StatusCode
	if response.RawResponse.ContentLength < 1024*64 { //64k
		respLogFields["raw"] = response.String()
	}
	log.InfoFieldsWithTrace(ro.Context, logrus.Fields{
		"request":  reqLogFields,
		"response": respLogFields,
		"elapsed":  elapsed.Seconds(),
	}, "requestsx request")

	return response, err
}
