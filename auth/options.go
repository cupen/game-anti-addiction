package auth

import (
	"fmt"
	"net/http"
	"net/url"
)

type Option func(c *Client)

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient == nil {
			panic(fmt.Errorf("nil http client"))
			// httpClient = http.DefaultClient
		}
		c.httpClient = httpClient
	}
}

func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

func WithDebugArgs(args url.Values) Option {
	return func(c *Client) {
		c.debugArgs = args
	}
}
