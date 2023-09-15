package user

import (
	"cert-gateway/cert/internal/errs"
	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/internal/syncx"
	"cert-gateway/cert/internal/types"
	"cert-gateway/cert/pkg/cache"
	"cert-gateway/utils"
	"context"
	"github.com/go-jose/go-jose/v3/json"
	"log"

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

	certsStr, err := json.Marshal(cache.CertCache.CaptureChanged(0))
	if err != nil {
		log.Panicln(err)
	}
	err = syncx.GlobalPubSub.Publish(string(certsStr))

	if err != nil {
		log.Panicln(err)
	}

	if req.Name == user.Name && req.Pass == user.Pass {
		return &types.UserLoginResp{
			Token: jwt,
		}, nil
	}

	return nil, errs.UnAuthorizationException
}
