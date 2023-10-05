package handler

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/go-resty/resty/v2"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"net/http"
	"testing"
)

func TestAcceptChallenge(t *testing.T) {
	type args struct {
		c *config.Config
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		{
			name: "run accept challenge and proxy success",
			args: struct{ c *config.Config }{c: &config.Config{Gateway: struct{ IssuerService string }{IssuerService: "http://127.0.0.1"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mux2 := &http.ServeMux{}
			mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			})
			server2 := &http.Server{
				Addr:    ":10086",
				Handler: mux2,
			}

			mux := &http.ServeMux{}
			handler := AcceptChallenge(tt.args.c)
			mux.HandleFunc("/", handler)

			server := &http.Server{
				Addr:    ":80",
				Handler: mux,
			}

			defer server.Shutdown(context.Background())
			defer server2.Shutdown(context.Background())

			go server2.ListenAndServe()
			go server.ListenAndServe()

			client := resty.New()

			resp, err := client.R().
				EnableTrace().
				Get("http://127.0.0.1:80/")
			if err != nil {
				t.Errorf("http request fail: [%s]", err)
				return
			}

			assert.Equal(t, resp.String(), "ok")

		})
	}
}

func TestHttpMiddleware(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "can create instance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := HttpMiddleware(func(w http.ResponseWriter, r *http.Request) {

			})
			if middleware == nil {
				t.Errorf("middleware is nil")
			}
		})
	}
}
