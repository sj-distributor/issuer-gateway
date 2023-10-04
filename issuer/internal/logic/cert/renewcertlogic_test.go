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

func TestNewRenewCertLogic(t *testing.T) {
	type args struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name string
		args args
		want *RenewCertLogic
	}{
		{
			name: "can create instance",
			args: struct {
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{ctx: context.Background(), svcCtx: &svc.ServiceContext{}},
			want: NewRenewCertLogic(context.Background(), &svc.ServiceContext{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRenewCertLogic(tt.args.ctx, tt.args.svcCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRenewCertLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewCertLogic_RenewCert(t *testing.T) {

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
			name: "can renew certificate",
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

			l := &RenewCertLogic{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}

			req := &types.CertificateRequest{
				Id: cert.Id,
			}

			_, err = l.RenewCert(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenewCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cert = &entity.Cert{}
			err = tt.fields.svcCtx.DB.Where("domain = ?", domain).Find(cert).Error
			if err != nil {
				t.Errorf("not found cert, domain: [%s]", err)
				return
			}

			if cert.Id == req.Id {
				t.Errorf("cert.Id must be not same: [%d], but: [%d]", cert.Id, req.Id)
				return
			}

			certificateDecrypt, privateKeyDecrypt, _, err := acme.DecryptCertificate(cert.Certificate, cert.PrivateKey, cert.IssuerCertificate, tt.fields.svcCtx.Config.Secret)

			if err != nil {
				t.Errorf("DecryptCertificate failed: [%s]", err)
				return
			}

			if certificateDecrypt == "" || privateKeyDecrypt == "" {
				t.Errorf("certificate can not be empty")
				t.Errorf("PrivateKey can not be empty")
				return
			}
		})
	}
}
