package runtime

type Action uint32

const (
	ActionContinue Action = 0
	ActionPause    Action = 1
)

type PeerType uint32

const (
	PeerTypeUnknown PeerType = 0
	PeerTypeLocal   PeerType = 1
	PeerTypeRemote  PeerType = 2
)

type LogLevel uint32

const (
	LogLevelTrace    LogLevel = 0
	LogLevelDebug    LogLevel = 1
	LogLevelInfo     LogLevel = 2
	LogLevelWarn     LogLevel = 3
	LogLevelError    LogLevel = 4
	LogLevelCritical LogLevel = 5
)

type Status uint32

const (
	StatusOk              Status = 0
	StatusNotFound        Status = 1
	StatusBadArgument     Status = 2
	StatusEmpty           Status = 7
	StatusCasMismatch     Status = 8
	StatusInternalFailure Status = 10
)

type MapType uint32

const (
	MapTypeHttpRequestHeaders       MapType = 0
	MapTypeHttpRequestTrailers      MapType = 1
	MapTypeHttpResponseHeaders      MapType = 2
	MapTypeHttpResponseTrailers     MapType = 3
	MapTypeHttpCallResponseHeaders  MapType = 7
	MapTypeHttpCallResponseTrailers MapType = 8
)

type BufferType uint32

const (
	BufferTypeHttpRequestBody      BufferType = 0
	BufferTypeHttpResponseBody     BufferType = 1
	BufferTypeDownstreamData       BufferType = 2
	BufferTypeUpstreamData         BufferType = 3
	BufferTypeHttpCallResponseBody BufferType = 4
)
