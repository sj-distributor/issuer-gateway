package cert

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"
	"github.com/pygzfei/issuer-gateway/pkg/acme"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCertFormUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCertFormUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCertFormUploadLogic {
	return &AddCertFormUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCertFormUploadLogic) AddCertFormUpload(req *types.AddCertFormUploadReq) (resp *types.AddOrRenewCertificateResp, err error) {

	certInfo := &certificate.Resource{
		PrivateKey:        []byte(req.PrivateKey),
		Certificate:       []byte(req.Certificate),
		IssuerCertificate: []byte(req.IssuerCertificate),
	}

	cert := &entity.Cert{
		Id: req.Id,
	}

	if l.svcCtx.DB.First(cert).Error != nil {
		return nil, err
	}

	certificateEncrypt, privateKeyEncrypt, issuerCertificateEncrypt, expire, err := acme.EncryptCertificate(certInfo, l.svcCtx.Config.Secret)
	if err != nil {
		return nil, err
	}

	cert.Expire = expire
	cert.Certificate = certificateEncrypt
	cert.PrivateKey = privateKeyEncrypt
	cert.IssuerCertificate = issuerCertificateEncrypt

	if l.svcCtx.DB.Save(cert).Error != nil {
		return nil, err
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

	return &types.AddOrRenewCertificateResp{}, nil
}
