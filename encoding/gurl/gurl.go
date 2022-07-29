package gurl

import (
	"github.com/Ravior/gserver/text/gstr"
	gconv2 "github.com/Ravior/gserver/util/gconv"
	"net/url"
	"strings"
)

const (
	fileUploadingKey = "@file:"
)

// Encode escapes the string so it can be safely placed
// inside a URL query.
func Encode(str string) string {
	return url.QueryEscape(str)
}

// Decode DecodeUint32 does the inverse transformation of Encode,
// converting each 3-byte encoded substring of the form "%AB" into the
// hex-decoded byte 0xAB.
// It returns an errors if any % is not followed by two hexadecimal
// digits.
func Decode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// URL-encode according to RFC 3986.
// See http://php.net/manual/en/function.rawurlencode.php.
func RawEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// DecodeUint32 URL-encoded strings.
// See http://php.net/manual/en/function.rawurldecode.php.
func RawDecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// Generate URL-encoded query string.
// See http://php.net/manual/en/function.http-build-query.php.
func BuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// Parse a URL and return its components.
// -1: all; 1: scheme; 2: host; 4: port; 8: user; 16: pass; 32: path; 64: query; 128: fragment.
// See http://php.net/manual/en/function.parse-url.php.
func ParseURL(str string, component int) (map[string]string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if component == -1 {
		component = 1 | 2 | 4 | 8 | 16 | 32 | 64 | 128
	}
	var components = make(map[string]string)
	if (component & 1) == 1 {
		components["scheme"] = u.Scheme
	}
	if (component & 2) == 2 {
		components["host"] = u.Hostname()
	}
	if (component & 4) == 4 {
		components["port"] = u.Port()
	}
	if (component & 8) == 8 {
		components["user"] = u.User.Username()
	}
	if (component & 16) == 16 {
		components["pass"], _ = u.User.Password()
	}
	if (component & 32) == 32 {
		components["path"] = u.Path
	}
	if (component & 64) == 64 {
		components["query"] = u.RawQuery
	}
	if (component & 128) == 128 {
		components["fragment"] = u.Fragment
	}
	return components, nil
}

// BuildParams builds the request string for the http client. The <params> can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter <noUrlEncode> specifies whether ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	// If given string/[]byte, converts and returns it directly as string.
	switch v := params.(type) {
	case string, []byte:
		return gconv2.String(params)
	case []interface{}:
		if len(v) > 0 {
			params = v[0]
		} else {
			params = nil
		}
	}
	// Else converts it to map and does the url encoding.
	m, urlEncode := gconv2.Map(params), true
	if len(m) == 0 {
		return gconv2.String(params)
	}
	if len(noUrlEncode) == 1 {
		urlEncode = !noUrlEncode[0]
	}
	// If there's file uploading, it ignores the url encoding.
	if urlEncode {
		for k, v := range m {
			if gstr.Contains(k, fileUploadingKey) || gstr.Contains(gconv2.String(v), fileUploadingKey) {
				urlEncode = false
				break
			}
		}
	}
	s := ""
	for k, v := range m {
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		s = gconv2.String(v)
		if urlEncode && len(s) > 6 && strings.Compare(s[0:6], fileUploadingKey) != 0 {
			s = Encode(s)
		}
		encodedParamStr += k + "=" + s
	}
	return
}
