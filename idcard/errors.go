package idcard

import "errors"

var (
	ErrInvalidLength = errors.New("invalid idnumbers.length, it should be 15 or 18")
)
