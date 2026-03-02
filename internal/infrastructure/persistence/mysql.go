package persistence

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxOpenConns = 100
	maxIdleConns = 20
)

// NewMySQL 创建 GORM 连接并配置连接池（高并发：MaxOpenConns=100, MaxIdleConns=20）
func NewMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("mysql open: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping: %w", err)
	}
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		log.Printf("warning: auto migrate users: %v", err)
	}
	return db, nil
}
