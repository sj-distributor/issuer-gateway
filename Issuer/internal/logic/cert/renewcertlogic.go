package cert

import (
	"cert-gateway/grpc/pb"
	"cert-gateway/issuer/internal/database/entity"
	"cert-gateway/issuer/internal/errs"
	"cert-gateway/pkg/acme"
	"context"
	"gorm.io/gorm"

	"cert-gateway/issuer/internal/svc"
	"cert-gateway/issuer/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RenewCertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenewCertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewCertLogic {
	return &RenewCertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenewCertLogic) RenewCert(req *types.CertificateRequest) (resp *types.AddOrRenewCertificateResp, err error) {

	cert := &entity.Cert{Id: req.Id}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查
		db := tx.First(cert)
		if db.Error != nil || db.RowsAffected == 0 {
			return errs.NotFoundException
		}

		// 2. renew
		certInfo, err := acme.ReqCertificate(l.svcCtx.Config.Env, cert.Email, cert.Domain)
		if err != nil {
			return err
		}

		certificateEncrypt, privateKeyEncrypt, issuerCertificateEncrypt, expire, err := acme.EncryptCertificate(certInfo, l.svcCtx.Config.Secret)
		if err != nil {
			return err
		}

		newCert := &entity.Cert{
			Domain:            cert.Domain,
			Email:             cert.Email,
			Target:            cert.Target,
			Certificate:       certificateEncrypt,
			PrivateKey:        privateKeyEncrypt,
			IssuerCertificate: issuerCertificateEncrypt,
			Expire:            expire,
		}

		// 3. 删
		db = tx.Delete(cert)
		if db.Error != nil || db.RowsAffected == 0 {
			return errs.NotFoundException
		}

		// 4. 新增 cert
		created := tx.Create(newCert)
		if created.Error != nil || created.RowsAffected == 0 {
			return errs.DatabaseError
		}

		err = l.svcCtx.SyncProvider.SyncCertificateToProvider(&pb.CertificateList{Certs: []*pb.Cert{
			{
				Id:                cert.Id,
				PrivateKey:        cert.PrivateKey,
				Certificate:       cert.Certificate,
				Domain:            cert.Domain,
				Target:            cert.Target,
				IssuerCertificate: cert.IssuerCertificate,
			},
		}})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.AddOrRenewCertificateResp{}, err
}
