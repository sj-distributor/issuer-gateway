package handler

import (
	"cert-gateway/gateway/internal/cache"
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
		proxy.Transport = &http.Transport{
			MaxIdleConns:        1024, // 最大空闲连接数
			MaxIdleConnsPerHost: 100,  // 每个主机的最大空闲连接数
			ForceAttemptHTTP2:   true,
			DisableKeepAlives:   false,
		}
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
