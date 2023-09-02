package entity

import "time"

type Cert struct {
	Id          uint64    `gorm:"type:bigint(20) UNSIGNED;not null;"`
	Domain      string    `gorm:"type:varchar(64);not null;index:idx_domain"`
	Email       string    `gorm:"type:varchar(64);not null;"`
	Certificate string    `gorm:"type:text;not null;" json:"certificate"`
	PrivateKey  string    `gorm:"type:text;not null;" json:"private_key"`
	Expire      time.Time `gorm:"type:datetime;not null;" json:"expire"`
	BaseEntity
}
