package handler

import (
	"cert-gateway/gateway/internal/cache"
	"crypto/tls"
	"errors"
	"fmt"
)

func CertificateInject(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, b := cache.GlobalCache.Get(info.ServerName); b {
		return &cert.TlS, nil
	}
	return &tls.Certificate{}, errors.New(fmt.Sprintf("%s: certificate not found", info.ServerName))
}
