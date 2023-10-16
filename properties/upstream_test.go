package properties

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetUpstreamAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamAddress, []byte("127.0.0.1"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamAddress()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", result)
}

func TestGetUpstreamPort(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamPort, serializeUint64(8080))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamPort()
	require.NoError(t, err)
	require.Equal(t, uint64(8080), result)
}

func TestGetUpstreamTlsVersion(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamTlsVersion, []byte("TLSv1.3"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamTlsVersion()
	require.NoError(t, err)
	require.Equal(t, "TLSv1.3", result)
}

func TestGetUpstreamSubjectLocalCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamSubjectLocalCertificate,
		[]byte("CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamSubjectLocalCertificate()
	require.NoError(t, err)
	require.Equal(t, "CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US", result)
}

func TestGetUpstreamSubjectPeerCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamSubjectPeerCertificate,
		[]byte("CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamSubjectPeerCertificate()
	require.NoError(t, err)
	require.Equal(t, "CN=example.com,OU=IT,O=example,L=San Francisco,ST=California,C=US", result)
}

func TestGetUpstreamDnsSanLocalCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamDnsSanLocalCertificate, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamDnsSanLocalCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetUpstreamDnsSanPeerCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamDnsSanPeerCertificate, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamDnsSanPeerCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetUpstreamUriSanLocalCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamUriSanLocalCertificate, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamUriSanLocalCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetUpstreamUriSanPeerCertificate(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamUriSanPeerCertificate, []byte("example.com"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamUriSanPeerCertificate()
	require.NoError(t, err)
	require.Equal(t, "example.com", result)
}

func TestGetUpstreamSha256PeerCertificateDigest(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamSha256PeerCertificateDigest,
		[]byte("b714f3d6f83efc2fddf80b8feda3e3b21b3e27b5"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamSha256PeerCertificateDigest()
	require.NoError(t, err)
	require.Equal(t, "b714f3d6f83efc2fddf80b8feda3e3b21b3e27b5", result)
}

func TestGetUpstreamLocalAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamLocalAddress, []byte("192.168.1.1"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamLocalAddress()
	require.NoError(t, err)
	require.Equal(t, "192.168.1.1", result)
}

func TestGetUpstreamTransportFailureReason(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(upstreamTransportFailureReason, []byte("connection closed"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetUpstreamTransportFailureReason()
	require.NoError(t, err)
	require.Equal(t, "connection closed", result)
}
