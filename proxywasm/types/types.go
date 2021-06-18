package types

import "errors"

type (
	Headers  = [][2]string
	Trailers = [][2]string
)

var (
	ErrorStatusNotFound    = errors.New("error status returned by host: not found")
	ErrorStatusBadArgument = errors.New("error status returned by host: bad argument")
	ErrorStatusEmpty       = errors.New("error status returned by host: empty")
	ErrorStatusCasMismatch = errors.New("error status returned by host: cas mismatch")
	ErrorInternalFailure   = errors.New("error status returned by host: internal failure")
)
