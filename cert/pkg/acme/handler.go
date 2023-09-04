package acme

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func AcceptChallenge() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		targetUrl := fmt.Sprintf("http://%s:5001%s", r.Host, r.RequestURI)
		target, _ := url.Parse(targetUrl)

		r.Host = target.Host
		r.URL = target

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}
	
}
