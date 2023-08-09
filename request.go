package telegraph

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// request define an API request
type request struct {
	method   string
	endpoint string
	secured  bool
	form     url.Values
	header   http.Header
	body     io.Reader
	fullURL  string
}

// setFormParam set param with key/value to request form body
func (r *request) setFormParam(key string, value interface{}) {
	if r.form == nil {
		r.form = url.Values{}
	}
	r.form.Set(key, fmt.Sprintf("%v", value))
}

func (r *request) validate() error {
	if r.form == nil {
		r.form = url.Values{}
	}
	return nil
}

// RequestOption define option type for request
type RequestOption func(*request)

func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *request) {
		if r.header == nil {
			r.header = http.Header{}
		}
		if replace {
			r.header.Set(key, value)
		} else {
			r.header.Add(key, value)
		}
	}
}
