package errors

import "errors"

var ErrDependencyFailure = errors.New("dependency failure")
var ErrInternalError = errors.New("internal error")
var ErrEntryNotFound = errors.New("entry not found")
