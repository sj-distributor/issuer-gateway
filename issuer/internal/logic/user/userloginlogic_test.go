package user

import (
	"context"
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/pygzfei/issuer-gateway/pkg/errs"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"testing"
)

func TestNewUserLoginLogic(t *testing.T) {
	type args struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name string
		args args
		want *UserLoginLogic
	}{
		{
			name: "can create instance",
			args: struct {
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{ctx: context.Background(), svcCtx: &svc.ServiceContext{}},
			want: NewUserLoginLogic(context.Background(), &svc.ServiceContext{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserLoginLogic(tt.args.ctx, tt.args.svcCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserLoginLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserLoginLogic_UserLogin(t *testing.T) {
	c := &config.Config{}

	c.Secret = "66d2e42661bc292f8237b4736a423a36"

	c.Issuer.User = struct {
		Name string
		Pass string
	}{Name: "123213", Pass: "2312sdkjhfksd"}

	type fields struct {
		Logger logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	type args struct {
		req *types.UserLoginReq
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *types.UserLoginResp
		wantErr  error
	}{
		{
			name: "run login success",
			fields: struct {
				Logger logx.Logger
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{
				Logger: logx.WithContext(context.Background()), ctx: context.Background(), svcCtx: &svc.ServiceContext{
					Config: *c,
				}},

			args:    struct{ req *types.UserLoginReq }{req: &types.UserLoginReq{Name: c.Issuer.User.Name, Pass: c.Issuer.User.Pass}},
			wantErr: nil,
		},
		{
			name: "run login fail",
			fields: struct {
				Logger logx.Logger
				ctx    context.Context
				svcCtx *svc.ServiceContext
			}{
				Logger: logx.WithContext(context.Background()), ctx: context.Background(), svcCtx: &svc.ServiceContext{
					Config: *c,
				}},
			args:    struct{ req *types.UserLoginReq }{req: &types.UserLoginReq{Name: "error", Pass: c.Issuer.User.Pass}},
			wantErr: errs.UnAuthorizationException,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &UserLoginLogic{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}

			gotResp, err := l.UserLogin(tt.args.req)

			if gotResp == nil && tt.wantErr == nil {
				t.Errorf("UserLogin() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			if (err != nil) && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("UserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
