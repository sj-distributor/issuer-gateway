package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

func HttpMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Sloww(r.RemoteAddr,
			logx.Field("Scheme", r.URL.Scheme),
			logx.Field("Method", r.Method),
			logx.Field("Host", r.Host),
			logx.Field("URL", r.URL),
			logx.Field("Header", r.Header),
		)
		next(w, r)
	}
}

// ReverseProxyHandler 根据证书配置的反向代理
func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	if cert, ok := cache.GlobalCache.Get(r.Host); ok {
		target, err := url.Parse(cert.Target)
		if err != nil {
			logx.Errorw("ReverseProxyHandler url.Parse(cert.Target)", logx.Field("cert.Target", cert.Target))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		r.Host = target.Host

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	} else {
		logx.Errorw("ReverseProxyHandler certificate not found", logx.Field("r.Host", r.Host))
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// ReverseProxyOrRedirect 启动http时, 判断是否强转https
func ReverseProxyOrRedirect(w http.ResponseWriter, r *http.Request) {
	cert, ok := cache.GlobalCache.Get(r.Host)

	target := r.URL

	if ok {
		if cert.Certificate != "" {
			httpsURL := "https://" + r.Host + r.URL.Path
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			return
		}

		target, err := url.Parse(cert.Target)
		if err != nil {
			logx.Errorw("ReverseProxyOrRedirect url.Parse(cert.Target)", logx.Field("cert.Target", cert.Target))
			http.NotFound(w, r)
			return
		}
		r.Host = target.Host
		r.URL = target
	}
	//if !ok {
	//	logx.Errorw("ReverseProxyOrRedirect certificate not found", logx.Field("r.Host", r.Host))
	//	http.NotFound(w, r)
	//	return
	//}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func AcceptChallenge(c *config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		targetUrl := fmt.Sprintf("%s%s", c.Gateway.IssuerAddr, r.RequestURI)

		target, err := url.Parse(targetUrl)
		if err != nil {
			logx.Errorw("AcceptChallenge err", logx.Field("error", err))
			http.NotFound(w, r)
			return
		}

		r.Host = target.Host
		r.URL = target

		fmt.Println()
		fmt.Println("r.URL :", r.URL)
		fmt.Println()

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}

}
