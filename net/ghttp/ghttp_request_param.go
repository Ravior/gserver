package ghttp

import (
	"bytes"
	"github.com/Ravior/gserver/core/internal/utils"
	"github.com/Ravior/gserver/encoding/gurl"
	"github.com/Ravior/gserver/errors/gcode"
	"github.com/Ravior/gserver/errors/gerror"
	"github.com/Ravior/gserver/internal/json"
	"github.com/Ravior/gserver/text/gregex"
	gstr2 "github.com/Ravior/gserver/text/gstr"
	"github.com/Ravior/gserver/util/gconv"
	"io/ioutil"
	"strings"
)

// parseQuery parses query string into r.queryMap.
func (r *Request) parseQuery() {
	if r.parsedQuery {
		return
	}
	r.parsedQuery = true
	if r.URL.RawQuery != "" {
		var err error
		r.queryMap, err = gstr2.Parse(r.URL.RawQuery)
		if err != nil {
			panic(gerror.WrapCode(gcode.CodeInvalidParameter, err, ""))
		}
	}
}

// parseBody parses the request raw data into r.rawMap.
// Note that it also supports JSON data from client request.
func (r *Request) parseBody() {
	if r.parsedBody {
		return
	}
	r.parsedBody = true
	// There's no data posted.
	if r.ContentLength == 0 {
		return
	}
	if body := r.GetBody(); len(body) > 0 {
		// Trim space/new line characters.
		body = bytes.TrimSpace(body)
		// JSON format checks.
		if body[0] == '{' && body[len(body)-1] == '}' {
			_ = json.UnmarshalUseNumber(body, &r.bodyMap)
		}
		// Default parameters decoding.
		if r.bodyMap == nil {
			r.bodyMap, _ = gstr2.Parse(r.GetBodyString())
		}
	}
}

// parseForm parses the request form for HTTP method PUT, POST, PATCH.
// The form data is pared into r.formMap.
//
// Note that if the form was parsed firstly, the request body would be cleared and empty.
func (r *Request) parseForm() {
	if r.parsedForm {
		return
	}
	r.parsedForm = true
	// There's no data posted.
	if r.ContentLength == 0 {
		return
	}
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		var err error
		if gstr2.Contains(contentType, "multipart/") {
			// multipart/form-data, multipart/mixed
			if err = r.ParseMultipartForm(1024 * 1024); err != nil {
				panic(gerror.WrapCode(gcode.CodeInvalidRequest, err, ""))
			}
		} else if gstr2.Contains(contentType, "form") {
			// application/x-www-form-urlencoded
			if err = r.Request.ParseForm(); err != nil {
				panic(gerror.WrapCode(gcode.CodeInvalidRequest, err, ""))
			}
		}
		if len(r.PostForm) > 0 {
			// Re-parse the form data using united parsing way.
			params := ""
			for name, values := range r.PostForm {
				// Invalid parameter name.
				// Only allow chars of: '\w', '[', ']', '-'.
				if !gregex.IsMatchString(`^[\w\-\[\]]+$`, name) && len(r.PostForm) == 1 {
					// It might be JSON/XML content.
					if s := gstr2.Trim(name + strings.Join(values, " ")); len(s) > 0 {
						if s[0] == '{' && s[len(s)-1] == '}' || s[0] == '<' && s[len(s)-1] == '>' {
							r.bodyContent = gconv.UnsafeStrToBytes(s)
							params = ""
							break
						}
					}
				}
				if len(values) == 1 {
					if len(params) > 0 {
						params += "&"
					}
					params += name + "=" + gurl.Encode(values[0])
				} else {
					if len(name) > 2 && name[len(name)-2:] == "[]" {
						name = name[:len(name)-2]
						for _, v := range values {
							if len(params) > 0 {
								params += "&"
							}
							params += name + "[]=" + gurl.Encode(v)
						}
					} else {
						if len(params) > 0 {
							params += "&"
						}
						params += name + "=" + gurl.Encode(values[len(values)-1])
					}
				}
			}
			if params != "" {
				if r.formMap, err = gstr2.Parse(params); err != nil {
					panic(gerror.WrapCode(gcode.CodeInvalidParameter, err, ""))
				}
			}
		}
	}
	// It parses the request body without checking the Content-Type.
	if r.formMap == nil {
		if r.Method != "GET" {
			r.parseBody()
		}
		if len(r.bodyMap) > 0 {
			r.formMap = r.bodyMap
		}
	}
}

// GetBody retrieves and returns request body content as bytes.
// It can be called multiple times retrieving the same body content.
func (r *Request) GetBody() []byte {
	if r.bodyContent == nil {
		r.bodyContent, _ = ioutil.ReadAll(r.Body)
		r.Body = utils.NewReadCloser(r.bodyContent, true)
	}
	return r.bodyContent
}

// GetBodyString retrieves and returns request body content as string.
// It can be called multiple times retrieving the same body content.
func (r *Request) GetBodyString() string {
	return gconv.UnsafeBytesToStr(r.GetBody())
}
