package cert

import (
	"cert-gateway/cert/internal/database"
	"cert-gateway/cert/internal/database/entity"
	"cert-gateway/cert/internal/database/hooks"
	"context"
	"fmt"
	"gorm.io/gorm"

	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCertsPagingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	DB     *gorm.DB
}

func NewGetCertsPagingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCertsPagingLogic {
	return &GetCertsPagingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		DB:     database.DB(),
	}
}

func (l *GetCertsPagingLogic) GetCertsPaging(req *types.GetCertsPagingReq) (resp *types.GetCertsPagingResp, err error) {

	var certs []entity.Cert

	db := l.DB.Scopes(hooks.Paging(req.Page, req.Size)).Order("id, expire")
	if req.Domain != "" {
		db = db.Where(fmt.Sprintf("domain like %q", "%"+req.Domain+"%"))
	}
	if req.Email != "" {
		db = db.Where(fmt.Sprintf("email like %q", "%"+req.Email+"%"))
	}
	err = db.Find(&certs).Error

	if err != nil {
		return nil, err
	}

	var certDtos []types.CertDto

	for _, cert := range certs {
		certDtos = append(certDtos, types.CertDto{
			Id:        cert.Id,
			Domain:    cert.Domain,
			Target:    cert.Target,
			Email:     cert.Email,
			Expire:    cert.Expire.Unix(),
			CreatedAt: cert.CreatedAt.Unix(),
		})
	}

	return &types.GetCertsPagingResp{
		Certs: certDtos,
	}, nil
}
