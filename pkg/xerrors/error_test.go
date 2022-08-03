package xerrors_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cadicallegari/user/pkg/xerrors"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := xerrors.New(xerrors.NotFound, "not_found", "user not found")
	assert.Equal(t, "type = NotFound code = not_found desc = user not found", err.Error())

	e, ok := xerrors.FromError(err)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.NotFound, e.Type)
		assert.Equal(t, "not_found", e.Code)
		assert.Equal(t, "user not found", e.Message)
	}
}

func TestNewf(t *testing.T) {
	err := xerrors.Newf(xerrors.NotFound, "not_found", "user %s not found", "x")
	assert.Equal(t, "type = NotFound code = not_found desc = user x not found", err.Error())

	e, ok := xerrors.FromError(err)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.NotFound, e.Type)
		assert.Equal(t, "not_found", e.Code)
		assert.Equal(t, "user x not found", e.Message)
	}
}

func TestNewf_WithUnknownError(t *testing.T) {
	origErr := errors.New("unknown error")

	err := xerrors.Newf(xerrors.Unavailable, "unavailable", "service unavailable: %w", origErr)
	assert.Equal(t, "type = Unavailable code = unavailable desc = service unavailable: unknown error", err.Error())
	assert.Equal(t, origErr, errors.Unwrap(err))

	e, ok := xerrors.FromError(err)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.Unavailable, e.Type)
		assert.Equal(t, "unavailable", e.Code)
		assert.Equal(t, "service unavailable: unknown error", e.Message)
	}
}

func TestNewf_WithXError(t *testing.T) {
	origErr := xerrors.New(xerrors.DeadlineExceeded, "deadline_exceeded", "service took more than 1sec to respond")
	assert.Equal(t, "type = DeadlineExceeded code = deadline_exceeded desc = service took more than 1sec to respond", origErr.Error())

	err := xerrors.Newf(xerrors.Unavailable, "unavailable", "service unavailable: %w", origErr)
	assert.Equal(t, "type = Unavailable code = unavailable desc = service unavailable: service took more than 1sec to respond", err.Error())
	assert.Equal(t, origErr, errors.Unwrap(err))

	e, ok := xerrors.FromError(err)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.Unavailable, e.Type)
		assert.Equal(t, "unavailable", e.Code)
		assert.Equal(t, "service unavailable: service took more than 1sec to respond", e.Message)
	}

	// Guarantee that after we wrap the original error still report the correct message
	assert.Equal(t, "type = DeadlineExceeded code = deadline_exceeded desc = service took more than 1sec to respond", origErr.Error())
}

func TestWrapError(t *testing.T) {
	origErr := xerrors.New(xerrors.NotFound, "not_found", "user not found")

	err := fmt.Errorf("wrap: %w", origErr)
	assert.Equal(t, "wrap: type = NotFound code = not_found desc = user not found", err.Error())

	e, ok := xerrors.FromError(err)
	if assert.True(t, ok) {
		assert.Equal(t, xerrors.NotFound, e.Type)
		assert.Equal(t, "not_found", e.Code)
		assert.Equal(t, "user not found", e.Message)
		assert.Equal(t, "type = NotFound code = not_found desc = user not found", e.Error())
	}
}

func TestFromContextError(t *testing.T) {
	testCases := []struct {
		Err  error
		Type xerrors.Type
		Code string
	}{
		{
			Err:  context.DeadlineExceeded,
			Type: xerrors.DeadlineExceeded,
			Code: "deadline_exceeded",
		},
		{
			Err:  context.Canceled,
			Type: xerrors.Canceled,
			Code: "canceled",
		},
		{
			Err:  errors.New("my embedded error"),
			Type: xerrors.OK,
		},
		{
			Err:  nil,
			Type: xerrors.OK,
		},
	}

	for _, tc := range testCases {
		err := xerrors.FromContextError(tc.Err)

		var e *xerrors.Error
		ok := errors.As(err, &e)
		if tc.Type == xerrors.OK {
			assert.Equal(t, tc.Err, err)
			continue
		}
		if assert.True(t, ok) {
			assert.Equal(t, tc.Type, e.Type)
			assert.Equal(t, tc.Code, e.Code)
		}
	}
}
