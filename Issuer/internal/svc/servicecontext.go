package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
	"issuer-gateway/issuer/internal/config"
	"issuer-gateway/issuer/internal/database"
	"issuer-gateway/issuer/internal/middleware"
	"issuer-gateway/issuer/internal/syncx"
	"issuer-gateway/pkg/driver"
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
