package xerrors_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cadicallegari/user/pkg/xerrors"
)

func TestInvalidFieldError_Is(t *testing.T) {
	var err error

	err = xerrors.NewInvalidFieldError("wrong_amount", "amount", "this amount is not expected")
	assert.True(t, errors.Is(err, xerrors.InvalidField))
	assert.True(t, errors.Is(err, &xerrors.Error{}))
	assert.False(t, errors.Is(err, xerrors.InvalidFieldType))
	assert.False(t, errors.Is(err, xerrors.MissingField))
	assert.False(t, errors.Is(err, xerrors.FieldOutOfRange))

	// Same type and code so it's a valid InvalidFieldType error.
	err = xerrors.NewInvalidFieldError("invalid_field_type", "amount", "this amount is not expected")
	assert.True(t, errors.Is(err, xerrors.InvalidFieldType))
	assert.True(t, errors.Is(err, xerrors.InvalidField))
	assert.True(t, errors.Is(err, &xerrors.Error{}))
	assert.False(t, errors.Is(err, xerrors.MissingField))
	assert.False(t, errors.Is(err, xerrors.FieldOutOfRange))

	err = xerrors.NewInvalidFieldTypeError("amount", reflect.TypeOf("").String(), reflect.TypeOf(1).String())
	assert.True(t, errors.Is(err, xerrors.InvalidFieldType))
	assert.True(t, errors.Is(err, xerrors.InvalidField))
	assert.True(t, errors.Is(err, &xerrors.Error{}))
	assert.False(t, errors.Is(err, xerrors.MissingField))
	assert.False(t, errors.Is(err, xerrors.FieldOutOfRange))

	err = xerrors.NewMissingFieldError("amount")
	assert.True(t, errors.Is(err, xerrors.MissingField))
	assert.True(t, errors.Is(err, xerrors.InvalidField))
	assert.True(t, errors.Is(err, &xerrors.Error{}))
	assert.False(t, errors.Is(err, xerrors.InvalidFieldType))
	assert.False(t, errors.Is(err, xerrors.FieldOutOfRange))

	err = xerrors.NewFieldOutOfRangeError("amount", 3, 8, 2)
	assert.True(t, errors.Is(err, xerrors.FieldOutOfRange))
	assert.True(t, errors.Is(err, xerrors.InvalidField))
	assert.True(t, errors.Is(err, &xerrors.Error{}))
	assert.False(t, errors.Is(err, xerrors.InvalidFieldType))
	assert.False(t, errors.Is(err, xerrors.MissingField))
}
