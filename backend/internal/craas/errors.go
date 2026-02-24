package craas

import "errors"

// ErrUnauthorized indicates that the request was not authorized (HTTP 401).
var ErrUnauthorized = errors.New("unauthorized")
