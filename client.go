package telegraph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	apiURL  = "https://api.telegra.ph/"
	baseURL = "https://telegra.ph/"
)

type (
	doFunc func(req *http.Request) (*http.Response, error)
	logger func(format string, v ...any)
)

func NewClient(accessToken string) *Client {
	return &Client{
		AccessToken: accessToken,
		BaseURL:     apiURL,
		UserAgent:   "Telegraph/golang",
		HTTPClient:  http.DefaultClient,
		Logger:      log.New(os.Stderr, "Telegraph-golang ", log.LstdFlags).Printf,
	}
}

// Client define API client
type Client struct {
	AccessToken string
	BaseURL     string
	UserAgent   string
	HTTPClient  *http.Client
	Debug       bool
	Logger      logger
	do          doFunc
}

func (c *Client) debug(format string, v ...any) {
	if c.Debug && c.Logger != nil {
		c.Logger(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) error {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}

	if err := r.validate(); err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)

	if r.secured {
		if c.AccessToken == "" {
			return ErrEmptyAccessToken
		}
		r.setFormParam("access_token", c.AccessToken)
	}

	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}

	if bodyString != "" {
		if header.Get("Content-Type") == "" {
			header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		r.body = bytes.NewBufferString(bodyString)
	}
	c.debug("full url: %s, body: %s", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	c.debug("method: %#+v, fullUrl: %#+v", r.method, r.fullURL)
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#+v", req)

	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	return data, nil
}
