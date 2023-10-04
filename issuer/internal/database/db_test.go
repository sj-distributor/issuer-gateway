package database

import (
	"fmt"
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestDB(t *testing.T) {
	c := config.Config{}
	c.Issuer.Mysql = struct {
		User string
		Pass string
		Host string
		Port string
		DB   string
	}{User: "root", Pass: "123456", Host: "127.0.0.1", Port: "3306", DB: "issuer-gateway"}

	Init(&c)

	tests := []struct {
		name string
		want *gorm.DB
	}{
		{
			name: "can get db",
			want: db,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := DB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		c *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "can init database",
			args: struct {
				c *config.Config
			}{
				c: func() *config.Config {
					c := config.Config{}
					c.Issuer.Mysql = struct {
						User string
						Pass string
						Host string
						Port string
						DB   string
					}{User: "root", Pass: "123456", Host: "127.0.0.1", Port: "3306", DB: "issuer-gateway"}
					return &c
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.c)
			s, err := DB().DB()
			if err != nil {
				t.Errorf("database init fail: [%s]", err)
				return
			}

			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s';", tt.args.c.Issuer.Mysql.DB, "cert")
			if err = s.QueryRow(query).Scan(&count); err != nil {
				t.Errorf("database init fail: [%s]", err)
				return
			}

			if count != 1 {
				t.Errorf("table not found")
			}

		})
	}
}
