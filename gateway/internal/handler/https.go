package handler

import (
	"crypto/tls"
	"net/http"
)

func Https() *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/", ReverseProxyHandler)

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	server.TLSConfig = &tls.Config{
		GetCertificate: CertificateInject,
	}

	return server
}
