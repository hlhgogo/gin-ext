// Package grequests implements a friendly API over Go's existing net/http library
package grequestsx

import (
	"context"
	"net/http"
	"testing"

	"github.com/levigross/grequests"
)

func TestDoRegularRequest(t *testing.T) {

	got, err := DoRegularRequest(http.MethodGet, "https://postman-echo.com/headers", &grequests.RequestOptions{
		Context: context.TODO(),
	}, Flags{EnableLog: true})

	if err != nil {
		t.Error(err)
		return
	}
	t.Log("got: ", got)
	var Res = struct {
		A map[string]string `json:"headers"`
	}{}
	err = got.JSON(&Res)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(Res)
}

func TestDoJsonRequest(t *testing.T) {
	var Res = struct {
		A map[string]string `json:"headers"`
	}{}

	err := JsonGet("https://postman-echo.com/headers", &grequests.RequestOptions{
		Context: context.TODO(),
	}, &Res, Flags{EnableLog: true})

	if err != nil {
		t.Error(err)
		return
	}
	t.Log("got: ", Res)
}
