package ghttp

import "github.com/Ravior/gserver/core/net/ghttp/internal/client"

type (
	Client         = client.Client
	ClientResponse = client.Response
)

// NewClient creates and returns a new HTTP client object.
func NewClient() *Client {
	return client.New()
}

// Get is a convenience method for sending GET request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Get(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("GET", url, data...)
}

// Put is a convenience method for sending PUT request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Put(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("PUT", url, data...)
}

// Post is a convenience method for sending POST request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Post(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("POST", url, data...)
}

// Delete is a convenience method for sending DELETE request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Delete(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("DELETE", url, data...)
}

// Head is a convenience method for sending HEAD request.
// NOTE that remembers CLOSING the response object when it'll never be used.
// Deprecated, please use g.Client().Head or NewClient().Head instead.
func Head(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("HEAD", url, data...)
}

// Patch is a convenience method for sending PATCH request.
// NOTE that remembers CLOSING the response object when it'll never be used.
// Deprecated, please use g.Client().Patch or NewClient().Patch instead.
func Patch(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("PATCH", url, data...)
}

// Connect is a convenience method for sending CONNECT request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Connect(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("CONNECT", url, data...)
}

// Options is a convenience method for sending OPTIONS request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Options(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("OPTIONS", url, data...)
}

// DoRequest is a convenience method for sending custom http method request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func DoRequest(method, url string, data ...interface{}) (*ClientResponse, error) {
	return client.New().DoRequest(method, url, data...)
}
