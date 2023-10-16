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
	return getPropertyString(upstreamAddress)
}

// GetUpstreamPort returns the upstream connection remote port.
func GetUpstreamPort() (uint64, error) {
	return getPropertyUint64(upstreamPort)
}

// GetUpstreamTlsVersion returns the TLS version of the upstream TLS connection.
func GetUpstreamTlsVersion() (string, error) {
	return getPropertyString(upstreamTlsVersion)
}

// GetUpstreamSubjectLocalCertificate returns the subject field of the local
// certificate in the upstream TLS connection.
func GetUpstreamSubjectLocalCertificate() (string, error) {
	return getPropertyString(upstreamSubjectLocalCertificate)
}

// GetUpstreamSubjectPeerCertificate returns the subject field of the peer
// certificate in the upstream TLS connection.
func GetUpstreamSubjectPeerCertificate() (string, error) {
	return getPropertyString(upstreamSubjectPeerCertificate)
}

// GetUpstreamDnsSanLocalCertificate returns the first DNS entry in the SAN
// field of the local certificate in the upstream TLS connection.
func GetUpstreamDnsSanLocalCertificate() (string, error) {
	return getPropertyString(upstreamDnsSanLocalCertificate)
}

// GetUpstreamDnsSanPeerCertificate returns the first DNS entry in the SAN
// field of the peer certificate in the upstream TLS connection.
func GetUpstreamDnsSanPeerCertificate() (string, error) {
	return getPropertyString(upstreamDnsSanPeerCertificate)
}

// GetUpstreamUriSanLocalCertificate returns the first URI entry in the SAN
// field of the local certificate in the upstream TLS connection.
func GetUpstreamUriSanLocalCertificate() (string, error) {
	return getPropertyString(upstreamUriSanLocalCertificate)
}

// GetUpstreamUriSanPeerCertificate returns the first URI entry in the SAN
// field of the peer certificate in the upstream TLS connection.
func GetUpstreamUriSanPeerCertificate() (string, error) {
	return getPropertyString(upstreamUriSanPeerCertificate)
}

// GetUpstreamSha256PeerCertificateDigest returns the SHA256 digest of the
// peer certificate in the upstream TLS connection if present.
func GetUpstreamSha256PeerCertificateDigest() (string, error) {
	return getPropertyString(upstreamSha256PeerCertificateDigest)
}

// GetUpstreamLocalAddress returns the local address of the upstream connection.
func GetUpstreamLocalAddress() (string, error) {
	return getPropertyString(upstreamLocalAddress)
}

// GetUpstreamTransportFailureReason returns the upstream transport failure
// reason e.g. certificate validation failed.
func GetUpstreamTransportFailureReason() (string, error) {
	return getPropertyString(upstreamTransportFailureReason)
}
