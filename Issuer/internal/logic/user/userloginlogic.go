package user

import (
	"cert-gateway/issuer/internal/errs"
	"cert-gateway/issuer/internal/svc"
	"cert-gateway/issuer/internal/types"
	"cert-gateway/utils"
	"context"
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
	user := l.svcCtx.Config.User
	jwt, err := utils.GenJwt(req.Pass, req.Name, l.svcCtx.Config.JWTSecret, l.svcCtx.Config.Secret)

	if req.Name == user.Name && req.Pass == user.Pass {
		return &types.UserLoginResp{
			Token: jwt,
		}, nil
	}

	return nil, errs.UnAuthorizationException
}
