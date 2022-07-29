package client

import (
	"bytes"
	"context"
	"github.com/Ravior/gserver/core/encoding/gurl"
	"github.com/Ravior/gserver/core/errors/gcode"
	"github.com/Ravior/gserver/core/errors/gerror"
	"github.com/Ravior/gserver/core/internal/json"
	"github.com/Ravior/gserver/core/os/gfile"
	"github.com/Ravior/gserver/core/text/gregex"
	"github.com/Ravior/gserver/core/text/gstr"
	"github.com/Ravior/gserver/core/util/gconv"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// DoRequest sends request with given HTTP method and data and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
//
// Note that it uses "multipart/form-data" as its Content-Type if it contains file uploading,
// else it uses "application/x-www-form-urlencoded". It also automatically detects the post
// content for JSON format, and for that it automatically sets the Content-Type as
// "application/json".
func (c *Client) DoRequest(method, url string, data ...interface{}) (resp *Response, err error) {
	req, err := c.prepareRequest(method, url, data...)
	if err != nil {
		return nil, err
	}
	resp, err = c.callRequest(req)
	return resp, err
}

// prepareRequest verifies request parameters, builds and returns http request.
func (c *Client) prepareRequest(method, url string, data ...interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	if len(c.prefix) > 0 {
		url = c.prefix + gstr.Trim(url)
	}
	var params string
	if len(data) > 0 {
		switch c.header["Content-Type"] {
		case "application/json":
			switch data[0].(type) {
			case string, []byte:
				params = gconv.String(data[0])
			default:
				if b, err := json.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}
		default:
			params = gurl.BuildParams(data[0])
		}
	}
	if method == "GET" {
		var bodyBuffer *bytes.Buffer
		if params != "" {
			switch c.header["Content-Type"] {
			case
				"application/json",
				"application/xml":
				bodyBuffer = bytes.NewBuffer([]byte(params))
			default:
				// It appends the parameters to the url
				// if http method is GET and Content-Type is not specified.
				if gstr.Contains(url, "?") {
					url = url + "&" + params
				} else {
					url = url + "?" + params
				}
				bodyBuffer = bytes.NewBuffer(nil)
			}
		} else {
			bodyBuffer = bytes.NewBuffer(nil)
		}
		if req, err = http.NewRequest(method, url, bodyBuffer); err != nil {
			return nil, err
		}
	} else {
		if strings.Contains(params, "@file:") {
			// File uploading request.
			var (
				buffer = bytes.NewBuffer(nil)
				writer = multipart.NewWriter(buffer)
			)
			for _, item := range strings.Split(params, "&") {
				array := strings.Split(item, "=")
				if len(array[1]) > 6 && strings.Compare(array[1][0:6], "@file:") == 0 {
					path := array[1][6:]
					if !gfile.Exists(path) {
						return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `"%s" does not exist`, path)
					}
					if file, err := writer.CreateFormFile(array[0], gfile.Basename(path)); err == nil {
						if f, err := os.Open(path); err == nil {
							if _, err = io.Copy(file, f); err != nil {
								if err := f.Close(); err != nil {
									//fmt.Errorf(c.ctx, `%+v`, err)
								}
								return nil, err
							}
							if err := f.Close(); err != nil {
								//fmt.Errorf(c.ctx, `%+v`, err)
							}
						} else {
							return nil, err
						}
					} else {
						return nil, err
					}
				} else {
					if err = writer.WriteField(array[0], array[1]); err != nil {
						return nil, err
					}
				}
			}
			// Close finishes the multipart message and writes the trailing
			// boundary end line to the output.
			if err = writer.Close(); err != nil {
				return nil, err
			}

			if req, err = http.NewRequest(method, url, buffer); err != nil {
				return nil, err
			} else {
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}
		} else {
			// Normal request.
			paramBytes := []byte(params)
			if req, err = http.NewRequest(method, url, bytes.NewReader(paramBytes)); err != nil {
				return nil, err
			} else {
				if v, ok := c.header["Content-Type"]; ok {
					// Custom Content-Type.
					req.Header.Set("Content-Type", v)
				} else if len(paramBytes) > 0 {
					if (paramBytes[0] == '[' || paramBytes[0] == '{') && json.Valid(paramBytes) {
						// Auto detecting and setting the post content format: JSON.
						req.Header.Set("Content-Type", "application/json")
					} else if gregex.IsMatchString(`^[\w\[\]]+=.+`, params) {
						// If the parameters passed like "name=value", it then uses form type.
						req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					}
				}
			}
		}
	}

	// Context.
	if c.ctx != nil {
		req = req.WithContext(c.ctx)
	} else {
		req = req.WithContext(context.Background())
	}
	// Custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}
	// It's necessary set the req.Host if you want to custom the host value of the request.
	// It uses the "Host" value from header if it's not empty.
	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	// Custom Cookie.
	if len(c.cookies) > 0 {
		headerCookie := ""
		for k, v := range c.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set("Cookie", headerCookie)
		}
	}
	// HTTP basic authentication.
	if len(c.authUser) > 0 {
		req.SetBasicAuth(c.authUser, c.authPass)
	}
	return req, nil
}

// callRequest sends request with give http.Request, and returns the responses object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) callRequest(req *http.Request) (resp *Response, err error) {
	resp = &Response{
		request: req,
	}

	for {
		if resp.Response, err = c.Do(req); err != nil {
			// The response might not be nil when err != nil.
			if resp.Response != nil {
				if err := resp.Response.Body.Close(); err != nil {
					//intlog.Errorf(c.ctx, `%+v`, err)
				}
			}
			if c.retryCount > 0 {
				c.retryCount--
				time.Sleep(c.retryInterval)
			} else {
				//return resp, err
				break
			}
		} else {
			break
		}
	}
	return resp, err
}
