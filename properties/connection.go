package properties

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
	result, err := getPropertyString(sourceAddress)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamRemotePort returns the remote port of the downstream connection.
func GetDownstreamRemotePort() (uint64, error) {
	result, err := getPropertyUint64(sourcePort)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetDownstreamLocalAddress returns the local address of the downstream connection.
func GetDownstreamLocalAddress() (string, error) {
	result, err := getPropertyString(destinationAddress)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamLocalPort returns the local port of the downstream connection.
func GetDownstreamLocalPort() (uint64, error) {
	result, err := getPropertyUint64(destinationPort)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetDownstreamConnectionID returns the connection ID of the downstream connection.
func GetDownstreamConnectionID() (uint64, error) {
	result, err := getPropertyUint64(connectionID)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// IsDownstreamConnectionTls returns true if the downstream connection is TLS.
func IsDownstreamConnectionTls() (bool, error) {
	result, err := getPropertyBool(connectionMtls)
	if err != nil {
		return false, err
	}
	return result, nil
}

// GetDownstreamRequestedServerName returns the requested server name of the
// downstream connection.
func GetDownstreamRequestedServerName() (string, error) {
	result, err := getPropertyString(connectionRequestedServerName)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamTlsVersion returns the TLS version of the downstream connection.
func GetDownstreamTlsVersion() (string, error) {
	result, err := getPropertyString(connectionTlsVersion)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamSubjectLocalCertificate returns the subject field of the local
// certificate in the downstream TLS connection.
func GetDownstreamSubjectLocalCertificate() (string, error) {
	result, err := getPropertyString(connectionSubjectLocalCert)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamSubjectPeerCertificate returns the subject field of the peer certificate
// in the downstream TLS connection.
func GetDownstreamSubjectPeerCertificate() (string, error) {
	result, err := getPropertyString(connectionSubjectPeerCert)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamDnsSanLocalCertificate returns The first DNS entry in the SAN field of
// the local certificate in the downstream TLS connection.
func GetDownstreamDnsSanLocalCertificate() (string, error) {
	result, err := getPropertyString(connectionDnsSanLocalCert)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamDnsSanPeerCertificate returns The first DNS entry in the SAN field of the
// peer certificate in the downstream TLS connection.
func GetDownstreamDnsSanPeerCertificate() (string, error) {
	result, err := getPropertyString(connectionDnsSanPeerCert)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamUriSanLocalCertificate returns the first URI entry in the SAN field of the
// local certificate in the downstream TLS connection
func GetDownstreamUriSanLocalCertificate() (string, error) {
	result, err := getPropertyString(connectionUriSanLocalCert)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamUriSanPeerCertificate returns The first URI entry in the SAN field of the
// peer certificate in the downstream TLS connection.
func GetDownstreamUriSanPeerCertificate() (string, error) {
	result, err := getPropertyString(connectionUriSanPeerCert)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamSha256PeerCertificateDigest returns the SHA256 digest of a peer certificate
// digest of the downstream connection.
func GetDownstreamSha256PeerCertificateDigest() (string, error) {
	result, err := getPropertyString(connectionSha256PeerCertDigest)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetDownstreamTerminationDetails returns the internal termination details of the connection
// (subject to change).
func GetDownstreamTerminationDetails() (string, error) {
	result, err := getPropertyString(connectionTerminationDetails)
	if err != nil {
		return "", err
	}
	return result, nil
}
