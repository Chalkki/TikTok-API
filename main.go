package main

import (
	"github.com/Chalkki/TikTok-API/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	controller.ConnectDb()
	r := gin.Default()
	initRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
