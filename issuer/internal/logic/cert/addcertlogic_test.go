package cert

import (
	"context"
	"fmt"
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/pygzfei/issuer-gateway/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
)

func TestAddCertLogic_AddCert(t *testing.T) {

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
			name: "run req certificate success",
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
				},
			},
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
			err = tt.fields.svcCtx.DB.First(cert).Error
			if err != nil {
				t.Errorf("not found cert, domain: [%s]", domain)
				return
			}

			l := &AddCertLogic{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}
			_, err = l.AddCert(&types.CertificateRequest{Id: cert.Id})
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewAddCertLogic(t *testing.T) {
	type args struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name string
		args args
		want *AddCertLogic
	}{
		{
			name: "can create instance",
			args: struct {
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{ctx: context.Background(), svcCtx: &svc.ServiceContext{}},
			want: NewAddCertLogic(context.Background(), &svc.ServiceContext{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddCertLogic(tt.args.ctx, tt.args.svcCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAddCertLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}
