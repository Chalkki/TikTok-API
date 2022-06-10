package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

//ConnectDb Connect to the mysql database
func ConnectDb() {
	var err error
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

//SetupDb see gorm documentations of gorm for AutoMigrate function details
// https://gorm.io/docs/migration.html#Auto-Migration
func SetupDb() {
	Db.AutoMigrate(&User{}, &UserLoginInfo{}, &Video{}, &Comment{}, &UserFavoriteInfo{}, &UserFollowInfo{})
}
