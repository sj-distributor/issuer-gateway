package handler

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"issuer-gateway/gateway/internal/cache"
	"issuer-gateway/gateway/internal/config"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxyHandler 根据证书配置的反向代理
func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	if cert, ok := cache.GlobalCache.Get(r.Host); ok {
		target, err := url.Parse(cert.Target)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		r.Host = target.Host

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// ReverseProxyOrRedirect 启动http时, 判断是否强转https
func ReverseProxyOrRedirect(w http.ResponseWriter, r *http.Request) {
	cert, ok := cache.GlobalCache.Get(r.Host)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if cert.Certificate != "" {
		httpsURL := "https://" + r.Host + r.URL.Path
		http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
		return
	}

	target, err := url.Parse(cert.Target)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	r.Host = target.Host

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func AcceptChallenge(c *config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		target, _ := url.Parse(c.Gateway.IssuerAddr)

		logx.Infov(target.Hostname())

		targetUrl := fmt.Sprintf("http://%s:5001%s", target.Hostname(), r.RequestURI)

		logx.Infof("Do challenge start: %s", targetUrl)

		r.Host = target.Host
		r.URL = target

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}

}
