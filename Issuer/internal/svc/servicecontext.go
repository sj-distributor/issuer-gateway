package svc

import (
	"cert-gateway/issuer/internal/config"
	"cert-gateway/issuer/internal/database"
	"cert-gateway/issuer/internal/middleware"
	"cert-gateway/issuer/internal/syncx"
	"cert-gateway/pkg/driver"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	Authorization rest.Middleware
	DB            *gorm.DB
	SyncProvider  driver.IProvider
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		Authorization: middleware.NewAuthorizationMiddleware(&c).Handle,
		DB:            database.DB(),
		SyncProvider:  syncx.Init(&c),
	}
}
