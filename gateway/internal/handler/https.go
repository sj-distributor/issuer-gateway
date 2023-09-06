package handler

import (
	"cert-gateway/gateway/internal/cache"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Https() *http.Server {

	//// 读取证书文件和私钥文件内容
	//certPEM, err := os.ReadFile("internal/cert/root.crt")
	//if err != nil {
	//	panic(err)
	//}
	//keyPEM, err := os.ReadFile("internal/cert/root.key")
	//if err != nil {
	//	panic(err)
	//}
	//
	//certificateEncrypt, privateKey, _, expire, err := acme.EncryptCertificate(&certificate.Resource{
	//	Certificate: certPEM,
	//	PrivateKey:  keyPEM,
	//}, "66d2e42661bc292f8237b4736a423a36")
	//
	//fmt.Println(certificateEncrypt, privateKey, expire.Unix())
	//
	//cache.GlobalCache.Set("test.anson.com", cache.Cert{
	//	Domain:      "test.anson.com",
	//	PrivateKey:  string(keyPEM),
	//	Certificate: string(certPEM),
	//	Target:      "http://192.167.167.167:9527",
	//})
	err := cache.GlobalCache.Sync()
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if cert, b := cache.GlobalCache.Get(r.Host); b {
			target, err := url.Parse(cert.Target)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			r.Host = target.Host

			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	server.TLSConfig = &tls.Config{
		GetCertificate: CertificateInject(),
	}

	return server
}
