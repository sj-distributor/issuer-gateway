package handler

import (
	"crypto/tls"
	"github.com/pygzfei/issuer-gateway/gateway/internal/cache"
	"github.com/pygzfei/issuer-gateway/gateway/internal/config"
	"github.com/pygzfei/issuer-gateway/utils"
	"reflect"
	"testing"
)

func TestCertificateInject(t *testing.T) {

	c := &config.Config{Secret: "66d2e42661bc292f8237b4736a423a36"}

	cache.Init(c)

	cache.GlobalCache.Set("test.anson.com", cache.Cert{
		Id:     uint64(utils.Id()),
		Domain: "test.anson.com",
		TlS:    tls.Certificate{},
	})

	type args struct {
		info *tls.ClientHelloInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *tls.Certificate
		wantErr bool
	}{
		{
			name:    "run get certificate success",
			args:    struct{ info *tls.ClientHelloInfo }{info: &tls.ClientHelloInfo{ServerName: "test.anson.com"}},
			want:    &tls.Certificate{},
			wantErr: false,
		},
		{
			name:    "run get certificate fail",
			args:    struct{ info *tls.ClientHelloInfo }{info: &tls.ClientHelloInfo{ServerName: "test1.anson.com"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, err := CertificateInject(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("CertificateInject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CertificateInject() got = %v, want %v", got, tt.want)
			}
		})
	}
}
