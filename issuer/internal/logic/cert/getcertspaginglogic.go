package cert

import (
	"context"
	"fmt"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/hooks"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"gorm.io/gorm"

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

	db := l.DB.Table("cert")
	if req.Domain != "" {
		db = db.Where(fmt.Sprintf("domain like %q", "%"+req.Domain+"%"))
	}
	if req.Email != "" {
		db = db.Where(fmt.Sprintf("email like %q", "%"+req.Email+"%"))
	}

	count := int64(0)
	err = db.Count(&count).Error
	if err != nil {
		return nil, err
	}

	err = db.Scopes(hooks.Paging(req.Page, req.Size)).Order("id, expire").Find(&certs).Error

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
			Expire:    cert.Expire,
			CreatedAt: cert.CreatedAt.Unix(),
		})
	}

	return &types.GetCertsPagingResp{
		Certs: certDtos,
		Total: uint64(count),
	}, nil
}
