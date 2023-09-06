package acme

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func AcceptChallenge() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		targetUrl := fmt.Sprintf("http://%s:5001%s", r.Host, r.RequestURI)
		target, _ := url.Parse(targetUrl)

		logx.Infof("Do challenge start: %s", targetUrl)
		
		r.Host = target.Host
		r.URL = target

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}

}
