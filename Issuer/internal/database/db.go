package database

import (
	"cert-gateway/issuer/internal/config"
	"cert-gateway/issuer/internal/database/entity"
	"cert-gateway/issuer/internal/database/hooks"
	"log"
	"time"

	"github.com/pygzfei/gorm-dbup/pkg"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func Init(c *config.Config) {
	database, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.Mysql.Dns, // DSN data source name
		DefaultStringSize:         256,         // string 类型字段的默认长度
		DisableDatetimePrecision:  true,        // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,        // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,        // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,       // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用数据库外键约束
	})

	// 插入时自动生成 雪花Id
	generateSnowflakeId := hooks.GenerateId()
	err = database.Callback().Create().Before("gorm:create").Register(generateSnowflakeId.Name, generateSnowflakeId.Initialize)

	db = database

	if c.Env == "dev" || c.Env == "debug" {
		err = database.AutoMigrate(&entity.Cert{})
		if err != nil {
			log.Fatalln(err)
		}
		db = database.Debug()
	} else if c.Env == "release" {
		err = database.Use(pkg.NewMigration(c.Mysql.Migration.Db, c.Mysql.Migration.Path))
		if err != nil {
			log.Fatalln(err)
		}
	}

	s, err := database.DB()
	if err != nil {
		log.Fatalln(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	s.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	s.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	s.SetConnMaxLifetime(time.Hour)
}

func DB() *gorm.DB {
	return db
}

func Close() {
	if sdb, e := db.DB(); e == nil {
		err := sdb.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
