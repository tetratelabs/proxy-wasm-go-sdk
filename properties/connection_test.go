package properties

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetDownstreamRemoteAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(sourceAddress, []byte("10.244.0.1:63649"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamRemoteAddress()
	require.NoError(t, err)
	require.Equal(t, "10.244.0.1:63649", result)
}

func TestGetDownstreamRemotePort(t *testing.T) {

	opt := proxytest.NewEmulatorOption().WithProperty(sourcePort, serializeUint64(63649))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamRemotePort()
	require.NoError(t, err)
	require.Equal(t, uint64(63649), result)
}

func TestGetDownstreamLocalAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(destinationAddress, []byte("10.244.0.13:80"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamLocalAddress()
	require.NoError(t, err)
	require.Equal(t, "10.244.0.13:80", result)
}

func TestGetDownstreamLocalPort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(destinationPort, serializeUint64(80))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamLocalPort()
	require.NoError(t, err)
	require.Equal(t, uint64(80), result)
}

func TestGetDownstreamConnectionID(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionID, serializeUint64(1771))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamConnectionID()
	require.NoError(t, err)
	require.Equal(t, uint64(1771), result)
}

func TestIsDownstreamConnectionTls(t *testing.T) {
	tests := []struct {
		name   string
		input  bool
		expect bool
	}{
		{
			name:   "Downstream Connection TLS: False",
			input:  false,
			expect: false,
		},
		{
			name:   "Downstream Connection TLS: True",
			input:  true,
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := proxytest.NewEmulatorOption().WithProperty(connectionMtls, serializeBool(tt.input))
			_, reset := proxytest.NewHostEmulator(opt)
			defer reset()

			result, err := IsDownstreamConnectionTls()
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}
func TestGetDownstreamRequestedServerName(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionRequestedServerName, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamRequestedServerName()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetDownstreamTlsVersion(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionTlsVersion, []byte("TLSv1.3"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamTlsVersion()
	require.NoError(t, err)
	require.Equal(t, "TLSv1.3", result)
}

func TestGetDownstreamSubjectLocalCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionSubjectLocalCert,
		[]byte("CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamSubjectLocalCertificate()
	require.NoError(t, err)
	require.Equal(t, "CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US", result)
}

func TestGetDownstreamSubjectPeerCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionSubjectPeerCert,
		[]byte("CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamSubjectPeerCertificate()
	require.NoError(t, err)
	require.Equal(t, "CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US", result)
}

func TestGetDownstreamDnsSanLocalCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionDnsSanLocalCert, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamDnsSanLocalCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetDownstreamDnsSanPeerCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionDnsSanPeerCert, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamDnsSanPeerCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetDownstreamUriSanLocalCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionUriSanLocalCert, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamUriSanLocalCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetDownstreamUriSanPeerCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionUriSanPeerCert, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamUriSanPeerCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetDownstreamSha256PeerCertificateDigest(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionSha256PeerCertDigest,
		[]byte("b714f3d6f83efc2fddf80b8feda3e3b21b3e27b5"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamSha256PeerCertificateDigest()
	require.NoError(t, err)
	require.Equal(t, "b714f3d6f83efc2fddf80b8feda3e3b21b3e27b5", result)
}

func TestGetDownstreamTerminationDetails(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionTerminationDetails,
		[]byte("connection closed"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamTerminationDetails()
	require.NoError(t, err)
	require.Equal(t, "connection closed", result)
}
