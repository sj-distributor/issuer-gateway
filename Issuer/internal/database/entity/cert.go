package entity

import "time"

type Cert struct {
	Id                uint64    `gorm:"type:bigint(20) UNSIGNED;not null;"`
	Domain            string    `gorm:"type:varchar(64);not null;uniqueIndex:unique_domain"`
	Target            string    `gorm:"type:varchar(255);not null;"`
	Email             string    `gorm:"type:varchar(64);not null;"`
	Certificate       string    `gorm:"type:text;" json:"certificate"`
	PrivateKey        string    `gorm:"type:text;" json:"private_key"`
	IssuerCertificate string    `gorm:"type:text;" json:"issuer_certificate"`
	Expire            time.Time `gorm:"type:datetime;index:idx_expire" json:"expire"`
	BaseEntity
}
