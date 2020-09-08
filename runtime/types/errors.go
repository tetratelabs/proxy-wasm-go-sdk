package types

import (
	"errors"
	"strconv"
)

var (
	ErrorStatusNotFound    = errors.New("not found")
	ErrorStatusBadArgument = errors.New("bad argument")
	ErrorStatusEmpty       = errors.New("empty")
	ErrorStatusCasMismatch = errors.New("cas mismatch")
	ErrorInternalFailure   = errors.New("internal failure")
)

func StatusToError(status Status) error {
	switch status {
	case StatusOK:
		return nil
	case StatusNotFound:
		return ErrorStatusNotFound
	case StatusBadArgument:
		return ErrorStatusBadArgument
	case StatusEmpty:
		return ErrorStatusEmpty
	case StatusCasMismatch:
		return ErrorStatusCasMismatch
	case StatusInternalFailure:
		return ErrorInternalFailure
	}
	return errors.New("unknown status code: " + strconv.Itoa(int(status)))
}
