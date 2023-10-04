package cert

import (
	"context"
	"fmt"
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"github.com/pygzfei/issuer-gateway/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
)

func TestAddCertFormUploadLogic_AddCertFormUpload(t *testing.T) {

	c := &config.Config{}

	c.Secret = "66d2e42661bc292f8237b4736a423a36"

	c.Issuer.Mysql = struct {
		User string
		Pass string
		Host string
		Port string
		DB   string
	}{User: "root", Pass: "123456", Host: "127.0.0.1", Port: "3306", DB: "issuer-gateway"}

	c.Logger = struct {
		Level    string
		Mode     string
		Path     string
		KeepDays int
		MaxSize  int
	}{Level: "debug", Mode: "console", Path: "", KeepDays: 0, MaxSize: 0}

	database.Init(c)

	type fields struct {
		Logger logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "can upload certificate",
			fields: struct {
				Logger logx.Logger
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{
				Logger: logx.WithContext(context.Background()), ctx: context.Background(), svcCtx: &svc.ServiceContext{
					DB:           database.DB(),
					Config:       *c,
					SyncProvider: &mockProvider{},
					AcmeProvider: &mockAcmeProvider{},
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			domainLogic := &AddDomainLogic{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}

			domain := fmt.Sprintf("%d.%d.com", utils.Id(), utils.Id())
			_, err := domainLogic.AddDomain(&types.AddDomainReq{Domain: domain, Email: "nsgzfei@gmail.com", Target: "https://anson.com"})

			cert := &entity.Cert{Domain: domain}
			err = tt.fields.svcCtx.DB.Where("domain = ?", domain).Find(cert).Error
			if err != nil {
				t.Errorf("not found cert, domain: [%s]", domain)
				return
			}

			l := &AddCertFormUploadLogic{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}

			req := &types.AddCertFormUploadReq{
				Id:                cert.Id,
				Certificate:       "-----BEGIN CERTIFICATE-----\nMIIELzCCAxegAwIBAgIUHOnNZ7Wsdxl8rRS5lT7rCAMPEy8wDQYJKoZIhvcNAQEL\nBQAwgaYxCzAJBgNVBAYTAkNOMRIwEAYDVQQIDAlHdWFuZ0RvbmcxEjAQBgNVBAcM\nCUd1YW5nWmhvdTEQMA4GA1UECgwHQW5zb24uRzETMBEGA1UECwwKQW5zb24gVGVh\nbTEXMBUGA1UEAwwOdGVzdC5hbnNvbi5jb20xLzAtBgkqhkiG9w0BCQEWIG5zZ3pm\nZWlAZ21haWwuY29tIHRlc3QuYW5zb24uY29tMB4XDTIzMDgyODAyMjU1MFoXDTIz\nMDkyNzAyMjU1MFowgaYxCzAJBgNVBAYTAkNOMRIwEAYDVQQIDAlHdWFuZ0Rvbmcx\nEjAQBgNVBAcMCUd1YW5nWmhvdTEQMA4GA1UECgwHQW5zb24uRzETMBEGA1UECwwK\nQW5zb24gVGVhbTEXMBUGA1UEAwwOdGVzdC5hbnNvbi5jb20xLzAtBgkqhkiG9w0B\nCQEWIG5zZ3pmZWlAZ21haWwuY29tIHRlc3QuYW5zb24uY29tMIIBIjANBgkqhkiG\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs2Xsy+UbbXEqTNPZP5l2XlYnOFhVfjA0YOuF\na8DY+6RUPlwmJQZpcu5wr0vOj40hi4MvDmxNddQ8OsTB+IS/id6bDU2L90LIrtj2\nGxAgxgQsIhSjf/2IiUC5sPOMnCtIGNcPLe8VURASeSL9cXpA8KxsPRsjLV9l/zUX\nWnZKd2/rmO2+T6GB98Iv1UkW8BPwlD/zoJwqcj4lHHhPXo0Wisa0fE0HGpjWXjP6\n8a459g3EiCcs93QuIs4UIxOjqbAWnCWfDSIsdHVEppvO3PenG3ntAxpjhdNhycKq\nm8gBvyZ8b4odih8i2vN1+eNmbkHXLsfk7l0GaqAJP/kCeXDWBQIDAQABo1MwUTAd\nBgNVHQ4EFgQUpg2x2H4JvKnapCUClr4OzlpB6dUwHwYDVR0jBBgwFoAUpg2x2H4J\nvKnapCUClr4OzlpB6dUwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC\nAQEAIE6V1dW5gQN4qe266z0+ssUWJjXkxVDdvGsv5dxdwU7FomP+43znnK5GKmV2\nzNBra0mWhw1sABh9XNhivFtjjmUckZ5bzxVjmOF3Uau1Fn0mDhIOy8zCzkyhj9z6\n6joMBx6TxgfNKFJ9ZP2dUXC9nDinmhSriO8zp8rTRCvGXukTwV91cezrpxztJcB8\n5yQbshUwR6PS6tA5wIOYmoGzz8TJa2ndKB6NZqT1gSo4+ZrpEuRScXpnUatsCT8Z\ndEqDuswRTXPB8eHP3gGMqfZ6Ats/OGsdHdTBE8sJ5F2ga348rLlOtyIPi1xKgGrE\n1Ib3GcazcTOu2ZOJ4WBP/zD2Cg==\n-----END CERTIFICATE-----",
				PrivateKey:        "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCzZezL5RttcSpM\n09k/mXZeVic4WFV+MDRg64VrwNj7pFQ+XCYlBmly7nCvS86PjSGLgy8ObE111Dw6\nxMH4hL+J3psNTYv3Qsiu2PYbECDGBCwiFKN//YiJQLmw84ycK0gY1w8t7xVREBJ5\nIv1xekDwrGw9GyMtX2X/NRdadkp3b+uY7b5PoYH3wi/VSRbwE/CUP/OgnCpyPiUc\neE9ejRaKxrR8TQcamNZeM/rxrjn2DcSIJyz3dC4izhQjE6OpsBacJZ8NIix0dUSm\nm87c96cbee0DGmOF02HJwqqbyAG/Jnxvih2KHyLa83X542ZuQdcux+TuXQZqoAk/\n+QJ5cNYFAgMBAAECggEALsw/4VB6vynuJux8l6KoxiMjSAeDBc/9WesWeu1rrPlJ\nIJtZN/9cMqcQrinQUJI4VfR6qgCGlF4w+AOrtfCrJoPzXp0EDhRV1YazbIvggMdF\n2/4WSKUSoPtJdWeTHooL3K79PrZHkUXoC8Gc66VAm4ffFHGn04Y3TUPEO8zv0Afq\nVffQjIMTQi2J2btKit/tdipp5aNwoN/uDnBzuAHOkb+zRWeLMp3xvpBeeuMkeks+\n7YJx0H2H0Uv6QPFdvD2GWMjQ16RhopUYkL46SMXDWF6nWBxr9c+1hAN4M1sPLOK7\n7OT9k+WRsa4WP+zosf5vf2jFJpjLZMlTzzBonfipAQKBgQDX0kX9PP2Lc1HXMhJ9\ndwbPADQgbZsrK8gNrhwA6Z2phmFyb74UOTA6LRfQ/DkSJfzOOgesqjA0v0s++Ya3\n53SVzs8GfHGPHFUK4JkdsdaS6Yyu0up60Ses7NLqgjqCxS4snmRvGMw8ObfZ4M05\nO5ig8kjpP13ZmS88uz1XjujC2QKBgQDUy8Y0P/z5MaYxec6sjC5ty8CX5n7OkwBZ\nlPTkT/S31dJQtREJP+i5eVBwOqyhM9vgAxyENgkH5r+iiReK4ZHWAENnkANA4PMi\nziyxyExh+CfyQdw7VIp5RQXZZNq5rSdCSSrLqL3LFqJk5xXO0QqzsWk7bGwJwTKd\n6S1a6v3ZDQKBgEJUf+o6ynoHcUnAO+qPoGoSV/L3fM8h35RExJqLMked42k2aqbw\nhJ/8p+s0+Z1YS4BeWWl5zOMJP+kU65Ct9CjurLYDnSssu/5h1O1JcPcqDHDWpfYl\nPhppltE4QR9b1rsj2x5B8tM3sgemjaxfYqNkk4AMV52+9MOnkEzOwT7RAoGACAbH\nqyDeweeFhUg663c+KRYOZaxkDBavZLGhqxr4+BYwoKqzwc2PUa+pwRH1gP8bxA9Z\n/AKtxIaHo/HX2X04qwHHiRh9huz1PtLYDLypZOifWRvy2qoNrxVTayfKuEY3vOBV\npOjwf8CSz4uH3w0zgiOm/H7SdGu9JQeulkdW2+ECgYAT+ylimK+uU1DwAP98nnw5\nF42Fz+1jPaW2//XezHvkWESfoe6RLJT7qmFYnwofjpdi8u4GVwKuOoO4pJZLjZm9\ne44yJdzPBMW0t2Ge7+Aqsy68Hrsg1qUangBdhoVQpeTE9LdGVu0EIJxx2qT2WfCF\nncpGXUztuYRIcK5Fz3HKrA==\n-----END PRIVATE KEY-----",
				IssuerCertificate: "",
			}
			_, err = l.AddCertFormUpload(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddCertFormUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cert = &entity.Cert{}
			err = tt.fields.svcCtx.DB.Where("domain = ?", domain).Find(cert).Error
			if err != nil {
				t.Errorf("not found cert, domain: [%s]", err)
				return
			}

			certificateDecrypt, privateKeyDecrypt, _, err := acme.DecryptCertificate(cert.Certificate, cert.PrivateKey, cert.IssuerCertificate, tt.fields.svcCtx.Config.Secret)

			if err != nil {
				t.Errorf("DecryptCertificate failed: [%s]", err)
				return
			}

			if certificateDecrypt != req.Certificate || privateKeyDecrypt != req.PrivateKey {
				t.Errorf("certificate is must be same, want: [%s], but: [%s] \n", certificateDecrypt, req.Certificate)
				t.Errorf("PrivateKey is must be same, want: [%s], but: [%s] \n", privateKeyDecrypt, req.PrivateKey)
				return
			}
		})
	}
}

func TestNewAddCertFormUploadLogic(t *testing.T) {
	type args struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name string
		args args
		want *AddCertFormUploadLogic
	}{
		{
			name: "can create instance",
			args: struct {
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{ctx: context.Background(), svcCtx: &svc.ServiceContext{}},
			want: NewAddCertFormUploadLogic(context.Background(), &svc.ServiceContext{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddCertFormUploadLogic(tt.args.ctx, tt.args.svcCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAddCertFormUploadLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}
