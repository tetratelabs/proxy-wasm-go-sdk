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
	return getPropertyString(sourceAddress)
}

// GetDownstreamRemotePort returns the remote port of the downstream connection.
func GetDownstreamRemotePort() (uint64, error) {
	return getPropertyUint64(sourcePort)
}

// GetDownstreamLocalAddress returns the local address of the downstream connection.
func GetDownstreamLocalAddress() (string, error) {
	return getPropertyString(destinationAddress)
}

// GetDownstreamLocalPort returns the local port of the downstream connection.
func GetDownstreamLocalPort() (uint64, error) {
	return getPropertyUint64(destinationPort)
}

// GetDownstreamConnectionID returns the connection ID of the downstream connection.
func GetDownstreamConnectionID() (uint64, error) {
	return getPropertyUint64(connectionID)
}

// IsDownstreamConnectionTls returns true if the downstream connection is TLS.
func IsDownstreamConnectionTls() (bool, error) {
	return getPropertyBool(connectionMtls)
}

// GetDownstreamRequestedServerName returns the requested server name of the
// downstream connection.
func GetDownstreamRequestedServerName() (string, error) {
	return getPropertyString(connectionRequestedServerName)
}

// GetDownstreamTlsVersion returns the TLS version of the downstream connection.
func GetDownstreamTlsVersion() (string, error) {
	return getPropertyString(connectionTlsVersion)
}

// GetDownstreamSubjectLocalCertificate returns the subject field of the local
// certificate in the downstream TLS connection.
func GetDownstreamSubjectLocalCertificate() (string, error) {
	return getPropertyString(connectionSubjectLocalCert)
}

// GetDownstreamSubjectPeerCertificate returns the subject field of the peer certificate
// in the downstream TLS connection.
func GetDownstreamSubjectPeerCertificate() (string, error) {
	return getPropertyString(connectionSubjectPeerCert)
}

// GetDownstreamDnsSanLocalCertificate returns The first DNS entry in the SAN field of
// the local certificate in the downstream TLS connection.
func GetDownstreamDnsSanLocalCertificate() (string, error) {
	return getPropertyString(connectionDnsSanLocalCert)
}

// GetDownstreamDnsSanPeerCertificate returns The first DNS entry in the SAN field of the
// peer certificate in the downstream TLS connection.
func GetDownstreamDnsSanPeerCertificate() (string, error) {
	return getPropertyString(connectionDnsSanPeerCert)
}

// GetDownstreamUriSanLocalCertificate returns the first URI entry in the SAN field of the
// local certificate in the downstream TLS connection
func GetDownstreamUriSanLocalCertificate() (string, error) {
	return getPropertyString(connectionUriSanLocalCert)
}

// GetDownstreamUriSanPeerCertificate returns The first URI entry in the SAN field of the
// peer certificate in the downstream TLS connection.
func GetDownstreamUriSanPeerCertificate() (string, error) {
	return getPropertyString(connectionUriSanPeerCert)
}

// GetDownstreamSha256PeerCertificateDigest returns the SHA256 digest of a peer certificate
// digest of the downstream connection.
func GetDownstreamSha256PeerCertificateDigest() (string, error) {
	return getPropertyString(connectionSha256PeerCertDigest)
}

// GetDownstreamTerminationDetails returns the internal termination details of the connection
// (subject to change).
func GetDownstreamTerminationDetails() (string, error) {
	return getPropertyString(connectionTerminationDetails)
}
