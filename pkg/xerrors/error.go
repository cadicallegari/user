package xerrors

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	err error
	// errIsXerror is used internally when generating Error message.
	// If .err is of type xerrors.Error we only get the message ignoring the .Type and .Code
	errIsXerror bool

	Type    Type   `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// xerr is used internally because other more specific errors wanted to
// embed Error, but we also have a .Error() method.
type xerr = Error

func new(t Type, code string, msg string) *Error {
	return &Error{
		Type:    t,
		Code:    code,
		Message: msg,
	}
}

func newf(t Type, code string, format string, a ...interface{}) *Error {
	tmpErr := fmt.Errorf(format, a...)
	if prevErr := errors.Unwrap(tmpErr); prevErr != nil {
		var e *Error
		if errors.As(prevErr, &e) {
			e.errIsXerror = true // avoid Error() return type and code
			tmpErr := fmt.Errorf(format, a...)
			e.errIsXerror = false

			err := new(t, code, tmpErr.Error())
			err.err = prevErr
			return err
		}

		// If we don't have an Error in the chain
		// we can call Error() to get the message and save to the Message field
		// we also can save the previous error
		err := new(t, code, tmpErr.Error())
		err.err = prevErr
		return err
	}
	return new(t, code, tmpErr.Error())
}

// New returns a Error representing t, code and msg.
func New(t Type, code string, msg string) error {
	return new(t, code, msg).Err()
}

// Newf returns New(t, code, fmt.Sprintf(format, a...)).
func Newf(t Type, code string, format string, a ...interface{}) error {
	return newf(t, code, format, a...).Err()
}

func FromError(err error) (*Error, bool) {
	if err == nil {
		return nil, true
	}
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return new(Unknown, "unknown_code", err.Error()), false
}

func Convert(err error) *Error {
	e, _ := FromError(err)
	return e
}

func (e *Error) Error() string {
	if e.errIsXerror {
		return e.Message
	}

	var b strings.Builder
	b.WriteString("type = ")
	b.WriteString(e.Type.String())
	b.WriteByte(' ')
	if e.Code != "" {
		b.WriteString("code = ")
		b.WriteString(e.Code)
		b.WriteByte(' ')
	}
	b.WriteString("desc = ")
	b.WriteString(e.Message)
	return b.String()
}

func (e *Error) Err() error {
	if e.Type == OK {
		return nil
	}
	return e
}

func (e *Error) Is(target error) bool {
	e2, ok := target.(*Error)
	// when target is just &Error{} type will be OK and we don't need to
	// check the Type and Code
	if ok && e2.Type != OK {
		return e.Type == e2.Type && e.Code == e2.Code
	}
	return ok
}

func (e *Error) Unwrap() error {
	return e.err
}

func ErrorType(err error) Type {
	// Don't use FromError to avoid allocation of OK status.
	if err == nil {
		return OK
	}
	if e, ok := FromError(err); ok {
		return e.Type
	}
	return Unknown
}

// FromContextError converts a context error into a Status.  It returns
// nil if err is nil and the original error if err is non-nil and not
// a context error.
func FromContextError(err error) error {
	switch err {
	case nil:
		return nil
	case context.DeadlineExceeded:
		return New(DeadlineExceeded, "deadline_exceeded", err.Error())
	case context.Canceled:
		return New(Canceled, "canceled", err.Error())
	default:
		return err
	}
}
