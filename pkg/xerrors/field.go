package xerrors

type InvalidFieldError struct {
	*xerr
	Field string `json:"field,omitempty"`
}

var (
	InvalidField     error = &InvalidFieldError{}
	MissingField     error = &InvalidFieldError{xerr: &Error{Type: InvalidArgument, Code: missingFieldCode}}
	InvalidFieldType error = &InvalidFieldError{xerr: &Error{Type: InvalidArgument, Code: invalidFieldTypeCode}}
	FieldOutOfRange  error = &FieldOutOfRangeError{}
)

func NewInvalidFieldError(code, field, msg string) *InvalidFieldError {
	e := &InvalidFieldError{}
	e.xerr = newf(InvalidArgument, code, msg)
	e.Field = field
	return e
}

func NewInvalidFieldErrorf(code, field, format string, a ...interface{}) *InvalidFieldError {
	e := &InvalidFieldError{}
	e.xerr = newf(InvalidArgument, code, format, a...)
	e.Field = field
	return e
}

func (e *InvalidFieldError) Is(target error) bool {
	e2, ok := target.(*InvalidFieldError)
	if !ok {
		return false
	}
	if e2.xerr == nil {
		return true
	}
	return e.Type == e2.Type && e.Code == e2.Code
}

func (e *InvalidFieldError) Unwrap() error {
	return e.xerr
}

const invalidFieldTypeCode = "invalid_field_type"

func NewInvalidFieldTypeError(field, exp, recv string) *InvalidFieldError {
	return NewInvalidFieldErrorf(invalidFieldTypeCode, field, "field %q expect %s but got %s", field, exp, recv)
}

const missingFieldCode = "missing_field"

func NewMissingFieldError(field string) *InvalidFieldError {
	return NewInvalidFieldErrorf(missingFieldCode, field, "field %q was not found", field)
}

type FieldOutOfRangeError struct {
	*InvalidFieldError
	Min int `json:"min"`
	Max int `json:"max"`
}

func NewFieldOutOfRangeError(field string, min, max, recv int) *FieldOutOfRangeError {
	e := &FieldOutOfRangeError{}
	e.InvalidFieldError = NewInvalidFieldErrorf("field_out_of_range", field, "field %q expect min %d and max %d but got %d", field, min, max, recv)
	e.Min = min
	e.Max = max
	return e
}

func (e *FieldOutOfRangeError) Is(err error) bool {
	_, ok := err.(*FieldOutOfRangeError)
	return ok
}

func (e *FieldOutOfRangeError) Unwrap() error {
	return e.InvalidFieldError
}
