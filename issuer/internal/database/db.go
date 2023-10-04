package database

import (
	"database/sql"
	"fmt"
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/entity"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database/hooks"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func Init(c *config.Config) {
	dbConfig := c.Issuer.Mysql

	dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True",
		dbConfig.User,
		dbConfig.Pass,
		dbConfig.Host,
		dbConfig.Port,
	))
	if err != nil {
		log.Fatalf("mysql init failed: [%s]", err)
	}
	_, err = dbConn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbConfig.DB))
	if err != nil {
		log.Fatalf("mysql init failed: [%s]", err)
	}

	defer func(dbConn *sql.DB) {
		_ = dbConn.Close()
	}(dbConn)

	database, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
			dbConfig.User,
			dbConfig.Pass,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DB,
		), // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用数据库外键约束
	})

	if err != nil {
		log.Fatalf("mysql init failed: [%s]", err)
	}

	// Automatically generate snowflake Id upon insertion
	generateSnowflakeId := hooks.GenerateId()
	err = database.Callback().Create().Before("gorm:create").Register(generateSnowflakeId.Name, generateSnowflakeId.Initialize)
	if err != nil {
		log.Fatalf("mysql init failed: [%s]", err)
	}

	err = database.AutoMigrate(&entity.Cert{})
	if err != nil {
		log.Fatalln(err)
	}

	db = database

	s, err := db.DB()
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
