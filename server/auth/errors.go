package auth

import (
	"errors"
)

var ErrAuthenticate = errors.New("Unauthorized user.")
