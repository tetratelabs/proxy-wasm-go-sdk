package properties

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

// This file hosts helper functions to retrieve connection-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#connection-attributes

var (
	sourceAddress                  = []string{"source", "address"}
	sourcePort                     = []string{"source", "port"}
	destinationAddress             = []string{"destination", "address"}
	destinationPort                = []string{"destination", "port"}
	connectionID                   = []string{"connection", "id"}
	connectionMtls                 = []string{"connection", "mtls"}
	connectionRequestedServerName  = []string{"connection", "requested_server_name"}
	connectionTlsVersion           = []string{"connection", "tls_version"}
	connectionSubjectLocalCert     = []string{"connection", "subject_local_certificate"}
	connectionSubjectPeerCert      = []string{"connection", "subject_peer_certificate"}
	connectionDnsSanLocalCert      = []string{"connection", "dns_san_local_certificate"}
	connectionDnsSanPeerCert       = []string{"connection", "dns_san_peer_certificate"}
	connectionUriSanLocalCert      = []string{"connection", "uri_san_local_certificate"}
	connectionUriSanPeerCert       = []string{"connection", "uri_san_peer_certificate"}
	connectionSha256PeerCertDigest = []string{"connection", "sha256_peer_certificate_digest"}
	connectionTerminationDetails   = []string{"connection", "termination_details"}
)

// GetDownstreamRemoteAddress returns the remote address of the downstream connection.
func GetDownstreamRemoteAddress() (string, error) {
	downstreamRemoteAddress, err := getPropertyString(sourceAddress)
	if err != nil {
		return "", err
	}
	return downstreamRemoteAddress, nil
}

// GetDownstreamRemotePort returns the remote port of the downstream connection.
func GetDownstreamRemotePort() (uint64, error) {
	downstreamRemotePort, err := getPropertyUint64(sourcePort)
	if err != nil {
		return 0, err
	}
	return downstreamRemotePort, nil
}

// GetDownstreamLocalAddress returns the local address of the downstream connection.
func GetDownstreamLocalAddress() (string, error) {
	downstreamLocalAddress, err := getPropertyString(destinationAddress)
	if err != nil {
		return "", err
	}
	return downstreamLocalAddress, nil
}

// GetDownstreamLocalPort returns the local port of the downstream connection.
func GetDownstreamLocalPort() (uint64, error) {
	downstreamLocalPort, err := getPropertyUint64(destinationPort)
	if err != nil {
		return 0, err
	}
	return downstreamLocalPort, nil
}

// GetDownstreamConnectionID returns the connection ID of the downstream connection.
func GetDownstreamConnectionID() (uint64, error) {
	downstreamConnectionId, err := getPropertyUint64(connectionID)
	if err != nil {
		return 0, err
	}
	return downstreamConnectionId, nil
}

// IsDownstreamConnectionTls returns true if the downstream connection is TLS.
func IsDownstreamConnectionTls() (bool, error) {
	downstreamConnectionTls, err := getPropertyBool(connectionMtls)
	if err != nil {
		return false, err
	}
	return downstreamConnectionTls, nil
}

// GetDownstreamRequestedServerName returns the requested server name of the downstream connection.
func GetDownstreamRequestedServerName() (string, error) {
	downstreamRequestedServerName, err := getPropertyString(connectionRequestedServerName)
	if err != nil {
		return "", err
	}
	return downstreamRequestedServerName, nil
}

// GetDownstreamTlsVersion returns the TLS version of the downstream connection.
func GetDownstreamTlsVersion() (string, error) {
	downstreamTlsVersion, err := getPropertyString(connectionTlsVersion)
	if err != nil {
		return "", err
	}
	return downstreamTlsVersion, nil
}

// GetDownstreamSubjectLocalCertificate returns the subject field of the local certificate in the downstream TLS connection.
func GetDownstreamSubjectLocalCertificate() (string, error) {
	downstreamSubjectLocalCertificate, err := getPropertyString(connectionSubjectLocalCert)
	if err != nil {
		return "", err
	}
	return downstreamSubjectLocalCertificate, nil
}

// GetDownstreamSubjectPeerCertificate returns the subject field of the peer certificate in the downstream TLS connection.
func GetDownstreamSubjectPeerCertificate() (string, error) {
	downstreamSubjectPeerCertificate, err := getPropertyString(connectionSubjectPeerCert)
	if err != nil {
		return "", err
	}
	return downstreamSubjectPeerCertificate, nil
}

// GetDownstreamDnsSanLocalCertificate returns The first DNS entry in the SAN field of the local certificate in the downstream TLS connection.
func GetDownstreamDnsSanLocalCertificate() (string, error) {
	downstreamDnsSanLocalCertificate, err := getPropertyString(connectionDnsSanLocalCert)
	if err != nil {
		return "", err
	}
	return downstreamDnsSanLocalCertificate, nil
}

// GetDownstreamDnsSanPeerCertificate returns The first DNS entry in the SAN field of the peer certificate in the downstream TLS connection.
func GetDownstreamDnsSanPeerCertificate() (string, error) {
	downstreamDnsSanPeerCertificate, err := getPropertyString(connectionDnsSanPeerCert)
	if err != nil {
		return "", err
	}
	return downstreamDnsSanPeerCertificate, nil
}

// GetDownstreamUriSanLocalCertificate returns the first URI entry in the SAN field of the local certificate in the downstream TLS connection
func GetDownstreamUriSanLocalCertificate() (string, error) {
	downstreamUriSanLocalCertificate, err := getPropertyString(connectionUriSanLocalCert)
	if err != nil {
		return "", err
	}
	return downstreamUriSanLocalCertificate, nil
}

// GetDownstreamUriSanPeerCertificate returns The first URI entry in the SAN field of the peer certificate in the downstream TLS connection.
func GetDownstreamUriSanPeerCertificate() (string, error) {
	downstreamUriSanPeerCertificate, err := getPropertyString(connectionUriSanPeerCert)
	if err != nil {
		return "", err
	}
	return downstreamUriSanPeerCertificate, nil
}

// GetDownstreamSha256PeerCertificateDigest returns the SHA256 digest of a peer certificate digest of the downstream connection.
func GetDownstreamSha256PeerCertificateDigest() ([]byte, error) {
	return proxywasm.GetProperty(connectionSha256PeerCertDigest)
}

// GetDownstreamTerminationDetails returns the internal termination details of the connection (subject to change).
func GetDownstreamTerminationDetails() (string, error) {
	downstreamTerminationDetails, err := getPropertyString(connectionTerminationDetails)
	if err != nil {
		return "", err
	}
	return downstreamTerminationDetails, nil
}
