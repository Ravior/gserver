package gerror

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Ravior/gserver/errors/gcode"
	"io"
)

// Error is custom errors for additional features.
type Error struct {
	error error      // Wrapped errors.
	stack stack      // Stack array, which records the stack information when this errors is created or wrapped.
	text  string     // Error text, which is created by New* functions.
	code  gcode.Code // Error code if necessary.
}

// Error implements the interface of Error, it returns all the errors as string.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	errStr := err.text
	if errStr == "" && err.code != nil {
		errStr = err.code.Message()
	}
	if err.error != nil {
		if errStr != "" {
			errStr += ": "
		}
		errStr += err.error.Error()
	}
	return errStr
}

// Code returns the errors code.
// It returns CodeNil if it has no errors code.
func (err *Error) Code() gcode.Code {
	if err == nil {
		return gcode.CodeNil
	}
	return err.code
}

// Cause returns the root cause errors.
func (err *Error) Cause() error {
	if err == nil {
		return nil
	}
	loop := err
	for loop != nil {
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				// Internal Error struct.
				loop = e
			} else {
				return loop.error
			}
		} else {
			// return loop
			// To be compatible with Case of https://github.com/pkg/errors.
			return errors.New(loop.text)
		}
	}
	return nil
}

// Format formats the frame according to the fmt.Formatter interface.
//
// %v, %s   : Print all the errors string;
// %-v, %-s : Print current level errors string;
// %+s      : Print full stack errors list;
// %+v      : Print the errors string and full stack errors list;
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('-'):
			if err.text != "" {
				io.WriteString(s, err.text)
			} else {
				io.WriteString(s, err.Error())
			}
		case s.Flag('+'):
			if verb == 's' {
				io.WriteString(s, err.Stack())
			} else {
				io.WriteString(s, err.Error()+"\n"+err.Stack())
			}
		default:
			io.WriteString(s, err.Error())
		}
	}
}

// Stack returns the stack callers as string.
// It returns an empty string if the <err> does not support stacks.
func (err *Error) Stack() string {
	if err == nil {
		return ""
	}
	var (
		loop   = err
		index  = 1
		buffer = bytes.NewBuffer(nil)
	)
	for loop != nil {
		buffer.WriteString(fmt.Sprintf("%d. %-v\n", index, loop))
		index++
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				loop = e
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s\n", index, loop.error.Error()))
				index++
				break
			}
		} else {
			break
		}
	}
	return buffer.String()
}

// Current creates and returns the current level errors.
// It returns nil if current level errors is nil.
func (err *Error) Current() error {
	if err == nil {
		return nil
	}
	return &Error{
		error: nil,
		stack: err.stack,
		text:  err.text,
		code:  err.code,
	}
}

// Next returns the next level errors.
// It returns nil if current level errors or the next level errors is nil.
func (err *Error) Next() error {
	if err == nil {
		return nil
	}
	return err.error
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (err *Error) MarshalJSON() ([]byte, error) {
	return []byte(`"` + err.Error() + `"`), nil
}
