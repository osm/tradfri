package tradfri

import (
	"errors"
)

var (
	ErrBadRequest         = errors.New("bad request")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("identifier/psk is incorrect")
)
