package internal

import (
	"errors"
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type BufferType uint32

const (
	BufferTypeHttpRequestBody      BufferType = 0
	BufferTypeHttpResponseBody     BufferType = 1
	BufferTypeDownstreamData       BufferType = 2
	BufferTypeUpstreamData         BufferType = 3
	BufferTypeHttpCallResponseBody BufferType = 4
	BufferTypeGrpcReceiveBuffer    BufferType = 5
	BufferTypeVMConfiguration      BufferType = 6
	BufferTypePluginConfiguration  BufferType = 7
	BufferTypeCallData             BufferType = 8
)

type MapType uint32

const (
	MapTypeHttpRequestHeaders       MapType = 0
	MapTypeHttpRequestTrailers      MapType = 1
	MapTypeHttpResponseHeaders      MapType = 2
	MapTypeHttpResponseTrailers     MapType = 3
	MapTypeHttpCallResponseHeaders  MapType = 6
	MapTypeHttpCallResponseTrailers MapType = 7
)

type MetricType uint32

const (
	MetricTypeCounter   = 0
	MetricTypeGauge     = 1
	MetricTypeHistogram = 2
)

type StreamType uint32

const (
	StreamTypeRequest    StreamType = 0
	StreamTypeResponse   StreamType = 1
	StreamTypeDownstream StreamType = 2
	StreamTypeUpstream   StreamType = 3
)

type Status uint32

const (
	StatusOK              Status = 0
	StatusNotFound        Status = 1
	StatusBadArgument     Status = 2
	StatusEmpty           Status = 7
	StatusCasMismatch     Status = 8
	StatusInternalFailure Status = 10
	StatusUnimplemented   Status = 12
)

//go:inline
func StatusToError(status Status) error {
	switch Status(status) {
	case StatusOK:
		return nil
	case StatusNotFound:
		return types.ErrorStatusNotFound
	case StatusBadArgument:
		return types.ErrorStatusBadArgument
	case StatusEmpty:
		return types.ErrorStatusEmpty
	case StatusCasMismatch:
		return types.ErrorStatusCasMismatch
	case StatusInternalFailure:
		return types.ErrorInternalFailure
	case StatusUnimplemented:
		return types.ErrorUnimplemented
	}
	return errors.New("unknown status code: " + strconv.Itoa(int(status)))
}
