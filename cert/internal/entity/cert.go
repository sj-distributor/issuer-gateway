package entity

type Cert struct {
	Id          uint64 `gorm:"type:bigint(20) UNSIGNED;not null;"`
	Domain      string `gorm:"type:varchar(64);not null;index:idx_domain"`
	Certificate string `gorm:"type:text;not null;"`
	PrivateKey  string `gorm:"type:text;not null;"`
	Expire      uint64 `gorm:"type:bigint(20) UNSIGNED;not null;"`
	BaseEntity
}
