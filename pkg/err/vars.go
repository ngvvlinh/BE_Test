package err

import (
	"errors"
)

var (
	ErrSizeLimit = errors.New("Size limit exceeded")
	ErrNoEntry   = errors.New("Entry is missing")
)
