package yzrt

import (
	"io"
	ghttp "net/http"
	"strings"
)

// Http is the built-in http singleton, accessible in Yz as `http`.
var Http = &_httpBoc{}

type _httpBoc struct{}

// Get fetches the given URI and returns the response body as a String.
func (h *_httpBoc) Get(uri String) *Thunk[String] {
	return Go(func() String {
		resp, err := ghttp.Get(uri.GoString())
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return NewString(string(body))
	})
}

// Post sends a POST request with the given body and returns the response body.
func (h *_httpBoc) Post(uri String, body String) *Thunk[String] {
	return Go(func() String {
		resp, err := ghttp.Post(uri.GoString(), "application/json", strings.NewReader(body.GoString()))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return NewString(string(respBody))
	})
}
