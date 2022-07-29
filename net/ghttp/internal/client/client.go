package client

import (
	"context"
	"crypto/tls"
	"github.com/Ravior/gserver/errors/gcode"
	"github.com/Ravior/gserver/errors/gerror"
	"net/http"
	"time"
)

type Client struct {
	http.Client                     // Underlying HTTP Client.
	ctx           context.Context   // Context for each request.
	header        map[string]string // Custom header map.
	cookies       map[string]string // Custom cookie map.
	prefix        string            // Prefix for request.
	authUser      string            // HTTP basic authentication: user.
	authPass      string            // HTTP basic authentication: pass.
	retryCount    int               // Retry count when request fails.
	retryInterval time.Duration     // Retry interval when request fails.
}

const (
	defaultClientAgent = "GHttpClient:1.0"
)

func New() *Client {
	client := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server in default.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		header: make(map[string]string),
	}
	client.header["User-Agent"] = defaultClientAgent
	return client
}

// SetHeader sets a custom HTTP header pair for the client.
func (c *Client) SetHeader(key, value string) *Client {
	c.header[key] = value
	return c
}

// SetHeaderMap sets custom HTTP headers with map.
func (c *Client) SetHeaderMap(m map[string]string) *Client {
	for k, v := range m {
		c.header[k] = v
	}
	return c
}

// SetCookie sets a cookie pair for the client.
func (c *Client) SetCookie(key, value string) *Client {
	c.cookies[key] = value
	return c
}

// SetCookieMap sets cookie items with map.
func (c *Client) SetCookieMap(m map[string]string) *Client {
	for k, v := range m {
		c.cookies[k] = v
	}
	return c
}

// SetAgent sets the User-Agent header for client.
func (c *Client) SetAgent(agent string) *Client {
	c.header["User-Agent"] = agent
	return c
}

// SetPrefix sets the request server URL prefix.
func (c *Client) SetPrefix(prefix string) *Client {
	c.prefix = prefix
	return c
}

// SetContentType sets HTTP content type for the client.
func (c *Client) SetContentType(contentType string) *Client {
	c.header["Content-Type"] = contentType
	return c
}

// SetBasicAuth sets HTTP basic authentication information for the client.
func (c *Client) SetBasicAuth(user, pass string) *Client {
	c.authUser = user
	c.authPass = pass
	return c
}

// SetTLSConfig sets the TLS configuration of client.
func (c *Client) SetTLSConfig(tlsConfig *tls.Config) error {
	if v, ok := c.Transport.(*http.Transport); ok {
		v.TLSClientConfig = tlsConfig
		return nil
	}
	return gerror.NewCode(gcode.CodeInternalError, `cannot set TLSClientConfig for custom Transport of the client`)
}

// SetTimeout sets the request timeout for the client.
func (c *Client) SetTimeout(t time.Duration) *Client {
	c.Client.Timeout = t
	return c
}

// SetCtx sets context for each request of this client.
func (c *Client) SetCtx(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

// SetRetry sets retry count and interval.
func (c *Client) SetRetry(retryCount int, retryInterval time.Duration) *Client {
	c.retryCount = retryCount
	c.retryInterval = retryInterval
	return c
}
