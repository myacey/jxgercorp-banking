package cstmerr

import "errors"

var (
	ErrUnknown      = errors.New("unknown error")
	ErrInvalidToken = errors.New("invalid auth token")
)
