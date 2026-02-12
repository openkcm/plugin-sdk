package pluginerrors

import "errors"

var (
	ErrServerRequired = errors.New("plugin server is required")
)
