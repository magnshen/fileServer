package main

import (
	"UploadServer/controler"
	"github.com/gin-gonic/gin"
)


func main() {
	e := gin.Default()
	e.MaxMultipartMemory = 4<<20
	e.POST("/fileServer/upload", controler.UploadHandle)
	e.GET("/fileServer/getProgress", controler.GetProgress)
	e.Run(":8899")
}
