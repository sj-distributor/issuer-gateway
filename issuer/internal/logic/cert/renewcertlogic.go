package cert

import (
	"context"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/errs"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"gorm.io/gorm"

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
		db := tx.First(cert)
		if db.Error != nil || db.RowsAffected == 0 {
			return errs.NotFoundException
		}

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

		db = tx.Delete(cert)
		if db.Error != nil || db.RowsAffected == 0 {
			return errs.NotFoundException
		}

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
