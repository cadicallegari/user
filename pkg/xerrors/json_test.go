package xerrors_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cadicallegari/user/pkg/xerrors"
)

func invalidJSON() error {
	var res interface{}
	return json.Unmarshal([]byte(`{"invalid"}`), &res)
}

func emptyJSON() error {
	var res interface{}
	return json.Unmarshal([]byte(``), &res)
}

func invalidTypeJSON() error {
	var res struct {
		Test bool `json:"test"`
	}
	return json.Unmarshal([]byte(`{"test": "invalid"}`), &res)
}

func TestFromJSONUnmarshal(t *testing.T) {
	var (
		err     error
		ok      bool
		jsonErr *xerrors.JSONUnmarshalError
	)

	err = xerrors.FromJSONUnmarshal(invalidJSON())
	ok = errors.As(err, &jsonErr)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.InvalidArgument, jsonErr.Type)
		assert.Equal(t, "invalid_json", jsonErr.Code)
		assert.Empty(t, jsonErr.Field)
		assert.EqualValues(t, 11, jsonErr.Offset)
	}

	err = xerrors.FromJSONUnmarshal(emptyJSON())
	ok = errors.As(err, &jsonErr)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.InvalidArgument, jsonErr.Type)
		assert.Equal(t, "invalid_json", jsonErr.Code)
		assert.Empty(t, jsonErr.Field)
		assert.EqualValues(t, 0, jsonErr.Offset)
	}

	err = xerrors.FromJSONUnmarshal(invalidTypeJSON())
	ok = errors.As(err, &jsonErr)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.InvalidArgument, jsonErr.Type)
		assert.Equal(t, "invalid_json_field", jsonErr.Code)
		assert.Equal(t, "test", jsonErr.Field)
		assert.EqualValues(t, 18, jsonErr.Offset)
	}

	prevErr := errors.New("non json error")
	err = xerrors.FromJSONUnmarshal(prevErr)
	ok = errors.As(err, &jsonErr)
	if assert.False(t, ok) {
		assert.Equal(t, prevErr, err)
	}

	err = xerrors.FromJSONUnmarshal(nil)
	ok = errors.As(err, &jsonErr)
	if assert.False(t, ok) {
		assert.Nil(t, err)
	}
}

func TestJSONUnmarshalError_Is(t *testing.T) {
	var err error

	err = xerrors.FromJSONUnmarshal(invalidJSON())
	assert.True(t, errors.Is(err, xerrors.JSONUnmarshal))
	assert.True(t, errors.Is(err, &xerrors.Error{}))

	err = xerrors.FromJSONUnmarshal(emptyJSON())
	assert.True(t, errors.Is(err, xerrors.JSONUnmarshal))
	assert.True(t, errors.Is(err, &xerrors.Error{}))

	err = xerrors.FromJSONUnmarshal(invalidTypeJSON())
	assert.True(t, errors.Is(err, xerrors.JSONUnmarshal))
	assert.True(t, errors.Is(err, &xerrors.Error{}))
}
