package handler

import (
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"testing"
)

func TestHttp(t *testing.T) {
	type args struct {
		c *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "http server can startup",
			args: struct{ c *config.Config }{c: &config.Config{Gateway: struct{ IssuerService string }{IssuerService: "https://anson.com"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Http(tt.args.c)
		})
	}
}
