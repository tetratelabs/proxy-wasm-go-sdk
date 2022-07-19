package log

import (
	"fmt"
	"strings"
)

// Enumeration of log levels.
//
// LevelDisabled is a special case to disable logging. It must remain last.
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelDisabled
)

// Level defines log levels.
type Level uint32

// ParseLevel converts a level string into a Level value.
// It returns an error if the level string does not match a known value.
func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "trace":
		return LevelTrace, nil
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	case "critical":
		return LevelCritical, nil
	case "disabled":
		return LevelDisabled, nil
	default:
		return 0, fmt.Errorf("unknown level string: %q", level)
	}
}

// String returns a string representation of the log level.
// It makes Level implement the fmt.Stringer interface.
func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	case LevelDisabled:
		return "disabled"
	default:
		panic("invalid log level")
	}
}
