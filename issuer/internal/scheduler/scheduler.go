package scheduler

import (
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/logic/cert"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"github.com/pygzfei/issuer-gateway/pkg/driver"
	"github.com/pygzfei/issuer-gateway/pkg/schedule"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strings"
	"time"
)

func NewScheduler(c *config.Config, syncProvider driver.IProvider, acmeProvider acme.IAcme) schedule.IScheduler {
	var scheduler schedule.IScheduler

	typeLower := strings.ToLower(c.Issuer.CheckExpireWithCron.Type)

	switch typeLower {
	case "redis":
		redis := c.Sync.Redis
		scheduler = schedule.NewRedisScheduler(redis.Addrs, redis.User, redis.Pass, redis.MasterName, redis.Db)
	case "cron":
		scheduler = schedule.NewCronScheduler()
	default:
		log.Fatalln("schedule type must be: 'Cron' or 'Redis' ")
		return nil
	}

	err := scheduler.StartAsync(c.Issuer.CheckExpireWithCron.Cron, func() {
		db := database.DB()

		var certs []entity.Cert
		err := db.Where("expire <= ?", time.Now().UTC().Unix()).Find(&certs).Order("id").Error
		if err != nil {
			logx.Errorf("Failed to scan expired certificates: [%s]", err)
			return
		}

		for _, certEntity := range certs {
			err := cert.Renew(c, db, syncProvider, acmeProvider, &certEntity)
			if err != nil {
				logx.Errorf("Failed to renew domain certificate: [%s], err: %s", certEntity.Domain, err)
				continue
			} else {
				logx.Infof("Success to renew domain certificate: [%s]", certEntity.Domain)
			}
		}
	})

	if err != nil {
		log.Fatalf("Failed to start scheduler: %s", err)
	}
	return scheduler
}
