package ghttp

import (
	"bytes"
	"net/http"
)

// Response is the ghttp response manager.
// Note that it implements the ghttp.ResponseWriter interface with buffering feature.
type Response struct {
	*ResponseWriter                 // Underlying ResponseWriter.
	Writer          *ResponseWriter // Alias of ResponseWriter.
	Request         *Request        // According request.
}

// NewResponse creates and returns a new Response object.
func NewResponse(w http.ResponseWriter) *Response {
	r := &Response{
		ResponseWriter: &ResponseWriter{
			writer: w,
			buffer: bytes.NewBuffer(nil),
		},
	}
	r.Writer = r.ResponseWriter
	return r
}
