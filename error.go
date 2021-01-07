package tradfri

import (
	"errors"
)

var (
	ErrBadRequest         = errors.New("bad request")
	ErrInvalidCredentials = errors.New("identifier/psk is incorrect")
)
