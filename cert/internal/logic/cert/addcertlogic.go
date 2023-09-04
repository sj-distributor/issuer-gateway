package cert

import (
	"cert-gateway/cert/internal/database/entity"
	"cert-gateway/cert/internal/errs"
	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/internal/types"
	"cert-gateway/cert/pkg/acme"
	"cert-gateway/cert/pkg/cache"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCertLogic {
	return &AddCertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCertLogic) AddCert(req *types.CertificateRequest) (resp *types.AddOrRenewCertificateResp, err error) {

	cert := &entity.Cert{
		Id: req.Id,
	}
	db := l.svcCtx.DB.First(cert)
	if db.Error != nil || db.RowsAffected == 0 {
		return nil, db.Error
	}

	certInfo, err := acme.ReqCertificate(l.svcCtx, cert.Email, cert.Domain)

	err = acme.EncryptCertificate(certInfo, cert, l.svcCtx.Config.Secret)

	if err != nil {
		return nil, err
	}

	tx := l.svcCtx.DB.Save(cert)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, errs.DatabaseError
	}

	// 更新缓存
	cache.CertCache.Set(cert.Id, types.Cert{
		PrivateKey:  cert.PrivateKey,
		Certificate: cert.Certificate,
		Domain:      cert.Domain,
	})

	return &types.AddOrRenewCertificateResp{}, nil
}
