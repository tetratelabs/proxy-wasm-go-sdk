package log_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/log"
)

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		level         string
		expectedLevel Level
		expectedErr   error
	}{
		{level: "trace", expectedLevel: LevelTrace},
		{level: "debug", expectedLevel: LevelDebug},
		{level: "info", expectedLevel: LevelInfo},
		{level: "warn", expectedLevel: LevelWarn},
		{level: "error", expectedLevel: LevelError},
		{level: "critical", expectedLevel: LevelCritical},
		{level: "disabled", expectedLevel: LevelDisabled},
		{level: "invalid", expectedErr: errors.New(`unknown level string: "invalid"`)},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			lvl, err := ParseLevel(tc.level)

			assert.Equal(t, tc.expectedLevel, lvl)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestLevel_String(t *testing.T) {
	testCases := []struct {
		level        Level
		exptectedStr string
	}{
		{level: LevelTrace, exptectedStr: "trace"},
		{level: LevelDebug, exptectedStr: "debug"},
		{level: LevelInfo, exptectedStr: "info"},
		{level: LevelWarn, exptectedStr: "warn"},
		{level: LevelError, exptectedStr: "error"},
		{level: LevelCritical, exptectedStr: "critical"},
		{level: LevelDisabled, exptectedStr: "disabled"},
	}

	for _, tc := range testCases {
		t.Run(tc.exptectedStr, func(t *testing.T) {
			str := tc.level.String()

			assert.Equal(t, tc.exptectedStr, str)
		})
	}
}
