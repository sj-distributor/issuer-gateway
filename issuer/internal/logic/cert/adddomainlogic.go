package cert

import (
	"context"
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDomainLogic {
	return &AddDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDomainLogic) AddDomain(req *types.AddDomainReq) (resp *types.AddOrRenewCertificateResp, err error) {

	cert := &entity.Cert{
		Email:  req.Email,
		Domain: req.Domain,
		Target: req.Target,
	}

	db := l.svcCtx.DB.Create(&cert)

	if db.RowsAffected == 0 || db.Error != nil {
		return nil, db.Error
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
		return nil, err
	}

	return &types.AddOrRenewCertificateResp{}, nil
}
