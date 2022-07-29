package gerror

import (
	"fmt"
	"github.com/Ravior/gserver/errors/gcode"
)

// apiCode is the interface for Code feature.
type apiCode interface {
	Error() string
	Code() gcode.Code
}

// New creates and returns an errors which is formatted from given text.
func New(text string) error {
	return &Error{
		stack: callers(),
		text:  text,
		code:  gcode.CodeNil,
	}
}

// Newf returns an errors that formats as the given format and args.
func Newf(format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  gcode.CodeNil,
	}
}

// NewCode creates and returns an errors that has errors code and given text.
func NewCode(code gcode.Code, text ...string) error {
	errText := ""
	if len(text) > 0 {
		errText = text[0]
	}
	return &Error{
		stack: callers(),
		text:  errText,
		code:  code,
	}
}

// NewCodef returns an errors that has errors code and formats as the given format and args.
func NewCodef(code gcode.Code, format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
		code:  code,
	}
}

// WrapCode wraps error with code and text.
// It returns nil if given err is nil.
func WrapCode(code gcode.Code, err error, text ...string) error {
	if err == nil {
		return nil
	}
	errText := ""
	if len(text) > 0 {
		errText = text[0]
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  errText,
		code:  code,
	}
}

// Code returns the error code of current error.
// It returns CodeNil if it has no error code or it does not implements interface Code.
func Code(err error) gcode.Code {
	if err != nil {
		if e, ok := err.(apiCode); ok {
			return e.Code()
		}
	}
	return gcode.CodeNil
}
