package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB
var err error

func ConnectDb() {
	//Connect to the mysql database
	dsn := MysqlConnectionLink + "?charset=utf8mb4"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}
	SetupDb()
}

func SetupDb() {
	//see gorm documentations of gorm for AutoMigrate function details
	// https://gorm.io/docs/migration.html#Auto-Migration
	Db.AutoMigrate(&User{}, &UserLoginInfo{}, &Video{}, &Comment{}, &UserFavoriteInfo{}, &VideoCommentInfo{})
}
