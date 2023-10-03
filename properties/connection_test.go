package properties

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
)

func TestGetDownstreamRemoteAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(sourceAddress, []byte("1.2.3.4"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamRemoteAddress()
	require.NoError(t, err)
	require.Equal(t, "1.2.3.4", result)
}

func TestGetDownstreamRemotePort(t *testing.T) {
	var u64 [8]byte
	binary.LittleEndian.PutUint64(u64[:], 1234)
	opt := proxytest.NewEmulatorOption().WithProperty(sourcePort, u64[:])
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamRemotePort()
	require.NoError(t, err)
	require.Equal(t, uint64(1234), result)
}

func TestGetDownstreamLocalAddress(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(destinationAddress, []byte("1.2.3.4"))
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamLocalAddress()
	require.NoError(t, err)
	require.Equal(t, "1.2.3.4", result)
}

func TestGetDownstreamLocalPort(t *testing.T) {
	var u64 [8]byte
	binary.LittleEndian.PutUint64(u64[:], 1234)
	opt := proxytest.NewEmulatorOption().WithProperty(destinationPort, u64[:])
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamLocalPort()
	require.NoError(t, err)
	require.Equal(t, uint64(1234), result)
}

func TestGetDownstreamConnectionID(t *testing.T) {
	var u64 [8]byte
	binary.LittleEndian.PutUint64(u64[:], 1234)
	opt := proxytest.NewEmulatorOption().WithProperty(connectionID, u64[:])
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamConnectionID()
	require.NoError(t, err)
	require.Equal(t, uint64(1234), result)
}

func TestIsDownstreamConnectionTls(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithProperty(connectionMtls, []byte{1})
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := IsDownstreamConnectionTls()
	require.NoError(t, err)
	require.Equal(t, true, result)
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
		[]byte{0x01, 0x02, 0x03, 0x04})
	_, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	result, err := GetDownstreamSha256PeerCertificateDigest()
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, result)
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
