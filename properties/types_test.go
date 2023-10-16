package properties

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvoyTrafficDirectionString(t *testing.T) {
	tests := []struct {
		input    EnvoyTrafficDirection
		expected string
	}{
		{EnvoyTrafficDirectionUnspecified, "UNSPECIFIED"},
		{EnvoyTrafficDirectionInbound, "INBOUND"},
		{EnvoyTrafficDirectionOutbound, "OUTBOUND"},
		{EnvoyTrafficDirection(9999), "UNSPECIFIED"},
	}

	for _, test := range tests {
		result := test.input.String()
		require.Equal(t, test.expected, result)
	}
}

func TestIstioTrafficInterceptionModeString(t *testing.T) {
	tests := []struct {
		input    IstioTrafficInterceptionMode
		expected string
	}{
		{IstioTrafficInterceptionModeNone, "NONE"},
		{IstioTrafficInterceptionModeTproxy, "TPROXY"},
		{IstioTrafficInterceptionModeRedirect, "REDIRECT"},
		{IstioTrafficInterceptionMode(9999), "REDIRECT"},
	}

	for _, test := range tests {
		result := test.input.String()
		require.Equal(t, test.expected, result)
	}
}

func TestParseIstioTrafficInterceptionMode(t *testing.T) {
	tests := []struct {
		input    string
		expected IstioTrafficInterceptionMode
		err      error
	}{
		{"NONE", IstioTrafficInterceptionModeNone, nil},
		{"TPROXY", IstioTrafficInterceptionModeTproxy, nil},
		{"REDIRECT", IstioTrafficInterceptionModeRedirect, nil},
		{"INVALID", IstioTrafficInterceptionModeRedirect, fmt.Errorf("invalid IstioTrafficInterceptionMode: INVALID")},
	}

	for _, test := range tests {
		result, err := ParseIstioTrafficInterceptionMode(test.input)
		require.Equal(t, test.expected, result)

		if test.err != nil {
			require.EqualError(t, err, test.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}
