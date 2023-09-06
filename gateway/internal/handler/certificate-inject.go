package handler

import (
	"cert-gateway/gateway/internal/cache"
	"crypto/tls"
	"log"
)

func CertificateInject() func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		domain := info.ServerName
		if cert, b := cache.GlobalCache.Get(domain); b {

			pair, err := tls.X509KeyPair([]byte(cert.Certificate), []byte(cert.PrivateKey))
			if err != nil {
				log.Fatalln(err)
			}

			return &pair, nil
		}
		return &tls.Certificate{}, nil
	}
}
