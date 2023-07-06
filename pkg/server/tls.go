package server

import (
	"crypto/tls"

	log "github.com/sirupsen/logrus"
)

func NewTLSConfig(conf *ServerConfig) *tls.Config {
	// Sensible default
	var tlsMinVersion uint16 = tls.VersionTLS12

	switch conf.TLSMinVersion {
	case "1.0":
		log.Warn("Setting the minimum version of TLS to 1.0 - this is discouraged. Please use 1.2 or 1.3")
		tlsMinVersion = tls.VersionTLS10
	case "1.1":
		log.Warn("Setting the minimum version of TLS to 1.1 - this is discouraged. Please use 1.2 or 1.3")
		tlsMinVersion = tls.VersionTLS11
	case "1.2":
		tlsMinVersion = tls.VersionTLS12
	case "1.3":
		tlsMinVersion = tls.VersionTLS13
	}

	strictCiphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	}

	// nil uses the default cipher suite
	var cipherConfig []uint16 = nil

	// TLS 1.3 does not support configuring the Cipher suites
	if tlsMinVersion != tls.VersionTLS13 && conf.HttpTLSStrictCiphers {
		cipherConfig = strictCiphers
	}

	return &tls.Config{
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		MinVersion:               tlsMinVersion,
		CipherSuites:             cipherConfig,
	}

}
