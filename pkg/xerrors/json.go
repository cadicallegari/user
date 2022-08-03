package xerrors

import (
	"encoding/json"
	"io"
)

type JSONUnmarshalError struct {
	*xerr
	Field  string `json:"field,omitempty"`
	Offset int64  `json:"offset,omitempty"`
}

var JSONUnmarshal error = &JSONUnmarshalError{}

func FromJSONUnmarshal(err error) error {
	e := &JSONUnmarshalError{}
	if err == io.EOF {
		e.xerr = newf(InvalidArgument, "invalid_json", "%w", err)
		return e
	}
	if jsonErr, ok := err.(*json.SyntaxError); ok {
		e.xerr = newf(InvalidArgument, "invalid_json", "%w", err)
		e.Offset = jsonErr.Offset
		return e
	}
	if jsonErr, ok := err.(*json.InvalidUTF8Error); ok {
		e.xerr = newf(InvalidArgument, "invalid_json", "unable to unmarshal %s: %w", jsonErr.S, err)
		return e
	}
	if jsonErr, ok := err.(*json.InvalidUnmarshalError); ok {
		if jsonErr.Type == nil {
			e.xerr = newf(InvalidArgument, "invalid_json", "unable to unmarshal: %w", err)
		} else {
			e.xerr = newf(InvalidArgument, "invalid_json", "unable to unmarshal %s: %w", jsonErr.Type.String(), err)
		}
		return e
	}
	if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		e.xerr = newf(InvalidArgument, "invalid_json_field", "unable to unmarshal %s type %s into value of type %s", jsonErr.Field, jsonErr.Value, jsonErr.Type.String())
		e.Field = jsonErr.Field
		e.Offset = jsonErr.Offset
		return e
	}
	return err
}

func (e *JSONUnmarshalError) Is(err error) bool {
	_, ok := err.(*JSONUnmarshalError)
	return ok
}

func (e *JSONUnmarshalError) Unwrap() error {
	return e.xerr
}
