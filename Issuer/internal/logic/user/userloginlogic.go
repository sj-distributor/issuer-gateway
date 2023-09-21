package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"issuer-gateway/issuer/internal/errs"
	"issuer-gateway/issuer/internal/svc"
	"issuer-gateway/issuer/internal/types"
	"issuer-gateway/utils"
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
	jwt, err := utils.GenJwt(req.Pass, req.Name, l.svcCtx.Config.JWTSecret, l.svcCtx.Config.Secret)

	if req.Name == user.Name && req.Pass == user.Pass {
		return &types.UserLoginResp{
			Token: jwt,
		}, nil
	}

	return nil, errs.UnAuthorizationException
}
