package cert

import (
	"cert-gateway/cert/internal/database/entity"
	"cert-gateway/cert/internal/errs"
	"cert-gateway/cert/internal/svc"
	"cert-gateway/cert/internal/types"
	"cert-gateway/cert/pkg/cache"
	"cert-gateway/pkg/acme"
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

	certInfo, err := acme.ReqCertificate(l.svcCtx.Config.Env, cert.Email, cert.Domain)

	certificateEncrypt, privateKeyEncrypt, issuerCertificateEncrypt, expire, err := acme.EncryptCertificate(certInfo, l.svcCtx.Config.Secret)
	if err != nil {
		return nil, err
	}

	cert.IssuerCertificate = issuerCertificateEncrypt
	cert.Certificate = certificateEncrypt
	cert.PrivateKey = privateKeyEncrypt
	cert.Expire = expire

	tx := l.svcCtx.DB.Save(cert)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, errs.DatabaseError
	}

	// 更新缓存
	cache.CertCache.Set(cert.Id, types.Cert{
		PrivateKey:  cert.PrivateKey,
		Certificate: cert.Certificate,
		Domain:      cert.Domain,
		Target:      cert.Target,
	})

	return &types.AddOrRenewCertificateResp{}, nil
}
