package svc

import (
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/middleware"
	"github.com/pygzfei/issuer-gateway/issuer/internal/syncx"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"github.com/pygzfei/issuer-gateway/pkg/driver"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	Authorization rest.Middleware
	DB            *gorm.DB
	SyncProvider  driver.IProvider
	AcmeProvider  acme.IAcme
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		Authorization: middleware.NewAuthorizationMiddleware(&c).Handle,
		DB:            database.DB(),
		SyncProvider:  syncx.Init(&c),
		AcmeProvider:  &acme.AcmeProvider{},
	}
}
