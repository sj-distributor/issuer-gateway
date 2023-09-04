package svc

import (
	"cert-gateway/cert/internal/config"
	"cert-gateway/cert/internal/database"
	"cert-gateway/cert/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	Authorization rest.Middleware
	DB            *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		Authorization: middleware.NewAuthorizationMiddleware(&c).Handle,
		DB:            database.DB(),
	}
}
