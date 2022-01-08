package httpserver

import "errors"

var (
	ErrPageSize  = errors.New("invalid page size")
	ErrPageRange = errors.New("invalid page limit")
)
