package ghttp

import (
	"encoding/json"
	"github.com/Ravior/gserver/core/util/gconv"
	"net/http"
)

// Write writes <content> to the response buffer.
func (r *Response) Write(content ...interface{}) {
	if len(content) == 0 {
		return
	}
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	for _, v := range content {
		switch value := v.(type) {
		case []byte:
			r.buffer.Write(value)
		case string:
			r.buffer.WriteString(value)
		default:
			r.buffer.WriteString(gconv.String(v))
		}
	}
}

// WriteJson writes <content> to the response with JSON format.
func (r *Response) WriteJson(content interface{}) error {
	// If given string/[]byte, response it directly to client.
	switch content.(type) {
	case string, []byte:
		r.Header().Set("Content-Type", "application/json")
		r.Write(gconv.String(content))
		return nil
	}
	// Else use json.Marshal function to encode the parameter.
	if b, err := json.Marshal(content); err != nil {
		return err
	} else {
		r.Header().Set("Content-Type", "application/json")
		r.Write(b)
	}
	return nil
}
