package handler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"issuer-gateway/gateway/internal/cache"
)

func CertificateInject(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, b := cache.GlobalCache.Get(info.ServerName); b {
		return &cert.TlS, nil
	}
	return &tls.Certificate{}, errors.New(fmt.Sprintf("%s: certificate not found", info.ServerName))
}
