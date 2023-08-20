package db

import (
	"douyin/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "gorm:gorm@tcp(192.168.100.128:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	err = db.AutoMigrate(&model.User{}, &model.Comment{}, &model.Follow{}, &model.Video{}, &model.Favorite{}, &model.Message{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}
