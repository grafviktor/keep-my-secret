package constant

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicateRecord = errors.New("duplicate record")
	ErrDeleted         = errors.New("deleted")
	ErrNoUserID        = errors.New("no user ID")
	ErrBadArgument     = errors.New("bad argument")
)

const (
	APIStatusSuccess = "success"
	APIStatusFail    = "fail"
	APIStatusError   = "error"
)

const (
	APIMessageUnauthorized = "unauthorized"
	APIMessageBadRequest   = "bad request"
	APIMessageServerError  = "server error"
	APIMessageNotFound     = "not found"
)
