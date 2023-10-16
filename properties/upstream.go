package properties

// This file hosts helper functions to retrieve upstream-related properties as described in:
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#upstream-attributes

var (
	upstreamAddress                     = []string{"upstream", "address"}
	upstreamPort                        = []string{"upstream", "port"}
	upstreamTlsVersion                  = []string{"upstream", "tls_version"}
	upstreamSubjectLocalCertificate     = []string{"upstream", "subject_local_certificate"}
	upstreamSubjectPeerCertificate      = []string{"upstream", "subject_peer_certificate"}
	upstreamDnsSanLocalCertificate      = []string{"upstream", "dns_san_local_certificate"}
	upstreamDnsSanPeerCertificate       = []string{"upstream", "dns_san_peer_certificate"}
	upstreamUriSanLocalCertificate      = []string{"upstream", "uri_san_local_certificate"}
	upstreamUriSanPeerCertificate       = []string{"upstream", "uri_san_peer_certificate"}
	upstreamSha256PeerCertificateDigest = []string{"upstream", "sha256_peer_certificate_digest"}
	upstreamLocalAddress                = []string{"upstream", "local_address"}
	upstreamTransportFailureReason      = []string{"upstream", "transport_failure_reason"}
)

// GetUpstreamAddress returns the upstream connection remote address.
func GetUpstreamAddress() (string, error) {
	result, err := getPropertyString(upstreamAddress)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamPort returns the upstream connection remote port.
func GetUpstreamPort() (uint64, error) {
	result, err := getPropertyUint64(upstreamPort)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetUpstreamTlsVersion returns the TLS version of the upstream TLS connection.
func GetUpstreamTlsVersion() (string, error) {
	result, err := getPropertyString(upstreamTlsVersion)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamSubjectLocalCertificate returns the subject field of the local
// certificate in the upstream TLS connection.
func GetUpstreamSubjectLocalCertificate() (string, error) {
	result, err := getPropertyString(upstreamSubjectLocalCertificate)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamSubjectPeerCertificate returns the subject field of the peer
// certificate in the upstream TLS connection.
func GetUpstreamSubjectPeerCertificate() (string, error) {
	result, err := getPropertyString(upstreamSubjectPeerCertificate)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamDnsSanLocalCertificate returns the first DNS entry in the SAN
// field of the local certificate in the upstream TLS connection.
func GetUpstreamDnsSanLocalCertificate() (string, error) {
	result, err := getPropertyString(upstreamDnsSanLocalCertificate)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamDnsSanPeerCertificate returns the first DNS entry in the SAN
// field of the peer certificate in the upstream TLS connection.
func GetUpstreamDnsSanPeerCertificate() (string, error) {
	result, err := getPropertyString(upstreamDnsSanPeerCertificate)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamUriSanLocalCertificate returns the first URI entry in the SAN
// field of the local certificate in the upstream TLS connection.
func GetUpstreamUriSanLocalCertificate() (string, error) {
	result, err := getPropertyString(upstreamUriSanLocalCertificate)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamUriSanPeerCertificate returns the first URI entry in the SAN
// field of the peer certificate in the upstream TLS connection.
func GetUpstreamUriSanPeerCertificate() (string, error) {
	result, err := getPropertyString(upstreamUriSanPeerCertificate)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamSha256PeerCertificateDigest returns the SHA256 digest of the
// peer certificate in the upstream TLS connection if present.
func GetUpstreamSha256PeerCertificateDigest() (string, error) {
	result, err := getPropertyString(upstreamSha256PeerCertificateDigest)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamLocalAddress returns the local address of the upstream connection.
func GetUpstreamLocalAddress() (string, error) {
	result, err := getPropertyString(upstreamLocalAddress)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetUpstreamTransportFailureReason returns the upstream transport failure
// reason e.g. certificate validation failed.
func GetUpstreamTransportFailureReason() (string, error) {
	result, err := getPropertyString(upstreamTransportFailureReason)
	if err != nil {
		return "", err
	}
	return result, nil
}
