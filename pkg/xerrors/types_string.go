// Code generated by "stringer -type=Type -output=types_string.go"; DO NOT EDIT.

package xerrors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OK-0]
	_ = x[Unknown-1]
	_ = x[Canceled-2]
	_ = x[DeadlineExceeded-3]
	_ = x[InvalidArgument-4]
	_ = x[NotFound-5]
	_ = x[AlreadyExists-6]
	_ = x[PermissionDenied-7]
	_ = x[OutOfRange-8]
	_ = x[FailedPrecondition-9]
	_ = x[Unimplemented-10]
	_ = x[Internal-11]
	_ = x[Unavailable-12]
	_ = x[Unauthenticated-13]
}

const _Type_name = "OKUnknownCanceledDeadlineExceededInvalidArgumentNotFoundAlreadyExistsPermissionDeniedOutOfRangeFailedPreconditionUnimplementedInternalUnavailableUnauthenticated"

var _Type_index = [...]uint8{0, 2, 9, 17, 33, 48, 56, 69, 85, 95, 113, 126, 134, 145, 160}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}