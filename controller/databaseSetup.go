package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB
var err error

func ConnectDb() {
	//Change the part "root:123456@tcp(127.0.0.1:3306)/MINITIKTOK"
	//to your mysql database.
	//Root is the mysql username, 123456 is the password, 3306 is port
	//MINITIKTOK is the name of the database.
	//It's recommended that using an empty database
	dsn := "root:123456@tcp(127.0.0.1:3306)/MINITIKTOK?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}
	SetupDb()
}

func SetupDb() {
	//see gorm documentations for AutoMigrate function details
	Db.AutoMigrate(&User{}, &UserLoginInfo{}, &Video{}, &Comment{})
}
