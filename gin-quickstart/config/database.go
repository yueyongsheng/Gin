package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库变量
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	dsn := "root:123456@tcp(localhost:3306)/grom?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Database connected successfully")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}