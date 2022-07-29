package client

import (
	"github.com/Ravior/gserver/core/util/gconv"
	"io/ioutil"
	"net/http"
)

// Response is the struct for client request response.
type Response struct {
	*http.Response
	request     *http.Request
	requestBody []byte
	cookies     map[string]string
}

// initCookie initializes the cookie map attribute of Response.
func (r *Response) initCookie() {
	if r.cookies == nil {
		r.cookies = make(map[string]string)
		// Response might be nil.
		if r != nil && r.Response != nil {
			for _, v := range r.Cookies() {
				r.cookies[v.Name] = v.Value
			}
		}
	}
}

// GetCookie retrieves and returns the cookie value of specified <key>.
func (r *Response) GetCookie(key string) string {
	r.initCookie()
	return r.cookies[key]
}

// GetCookieMap retrieves and returns a copy of current cookie values map.
func (r *Response) GetCookieMap() map[string]string {
	r.initCookie()
	m := make(map[string]string, len(r.cookies))
	for k, v := range r.cookies {
		m[k] = v
	}
	return m
}

// ReadAll retrieves and returns the response content as []byte.
func (r *Response) ReadAll() []byte {
	// Response might be nil.
	if r == nil || r.Response == nil {
		return []byte{}
	}
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return nil
	}
	return body
}

// ReadAllString retrieves and returns the response content as string.
func (r *Response) ReadAllString() string {
	return gconv.UnsafeBytesToStr(r.ReadAll())
}

// Close closes the response when it will never be used.
func (r *Response) Close() error {
	if r == nil || r.Response == nil || r.Response.Close {
		return nil
	}
	r.Response.Close = true
	return r.Response.Body.Close()
}
