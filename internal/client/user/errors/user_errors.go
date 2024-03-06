package errors

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrUserUnauthorized = errors.New("user unauthorized")
