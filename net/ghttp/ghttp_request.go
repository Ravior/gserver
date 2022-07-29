package ghttp

import (
	"fmt"
	"github.com/Ravior/gserver/text/gregex"
	"net/http"
	"strings"
)

// Request is the context object for a request.
type Request struct {
	*http.Request
	Response    *Response              // Corresponding Response of this request.
	parsedHost  string                 // The parsed host name for current host used by GetHost function.
	clientIp    string                 // The parsed client ip for current host used by GetClientIp function.
	parsedQuery bool                   // A bool marking whether the GET parameters parsed.
	queryMap    map[string]interface{} // Query parameters map, which is nil if there's no query string.
	parsedBody  bool                   // A bool marking whether the request body parsed.
	formMap     map[string]interface{} // Form parameters map, which is nil if there's no form data from client.
	bodyMap     map[string]interface{} // Body parameters map, which might be nil if there're no body content.
	bodyContent []byte                 // Request body content.
	paramsMap   map[string]interface{} // Custom parameters map.
	parsedForm  bool                   // A bool marking whether request Form parsed for HTTP method PUT, POST, PATCH.
}

// NewRequest creates and returns a new request object.
func NewRequest(r *http.Request, w http.ResponseWriter) *Request {
	request := &Request{
		Request:  r,
		Response: NewResponse(w),
	}
	request.Response.Request = request
	return request
}

// GetHost returns current request host name, which might be a domain or an IP without port.
func (r *Request) GetHost() string {
	if len(r.parsedHost) == 0 {
		array, _ := gregex.MatchString(`(.+):(\d+)`, r.Host)
		if len(array) > 1 {
			r.parsedHost = array[1]
		} else {
			r.parsedHost = r.Host
		}
	}
	return r.parsedHost
}

// GetClientIp returns the client ip of this request without port.
// Note that this ip address might be modified by client header.
func (r *Request) GetClientIp() string {
	if len(r.clientIp) == 0 {
		realIps := r.Header.Get("X-Forwarded-For")
		if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
			ipArray := strings.Split(realIps, ",")
			r.clientIp = ipArray[0]
		}
		if r.clientIp == "" || strings.EqualFold("unknown", realIps) {
			r.clientIp = r.Header.Get("Proxy-Client-IP")
		}
		if r.clientIp == "" || strings.EqualFold("unknown", realIps) {
			r.clientIp = r.Header.Get("WL-Proxy-Client-IP")
		}
		if r.clientIp == "" || strings.EqualFold("unknown", realIps) {
			r.clientIp = r.Header.Get("HTTP_CLIENT_IP")
		}
		if r.clientIp == "" || strings.EqualFold("unknown", realIps) {
			r.clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
		}
		if r.clientIp == "" || strings.EqualFold("unknown", realIps) {
			r.clientIp = r.Header.Get("X-Real-IP")
		}
		if r.clientIp == "" || strings.EqualFold("unknown", realIps) {
			r.clientIp = r.GetRemoteIp()
		}
	}
	return r.clientIp
}

// GetRemoteIp returns the ip from RemoteAddr.
func (r *Request) GetRemoteIp() string {
	array, _ := gregex.MatchString(`(.+):(\d+)`, r.RemoteAddr)
	if len(array) > 1 {
		return array[1]
	}
	return r.RemoteAddr
}

// GetUrl returns current URL of this request.
func (r *Request) GetUrl() string {
	scheme := "ghttp"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf(`%s://%s%s`, scheme, r.Host, r.URL.String())
}

// GetReferer returns referer of this request.
func (r *Request) GetReferer() string {
	return r.Header.Get("Referer")
}
