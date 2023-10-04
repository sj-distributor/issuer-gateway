package user

import (
	"context"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/pygzfei/issuer-gateway/pkg/errs"
	"github.com/pygzfei/issuer-gateway/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	user := l.svcCtx.Config.Issuer.User

	if req.Name == user.Name && req.Pass == user.Pass {
		jwt, _ := utils.GenJwt(req.Pass, req.Name, l.svcCtx.Config.JWTSecret, l.svcCtx.Config.Secret)
		return &types.UserLoginResp{
			Token: jwt,
		}, nil
	}

	return nil, errs.UnAuthorizationException
}
