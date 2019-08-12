package main

import (
	"./controler"
	"github.com/gin-gonic/gin"
)


func main() {
	e := gin.Default()
	e.POST("/fileServer/uploadAppend", controler.AppendHandle)
	e.POST("/fileServer/uploadNewFile", controler.UploadNewFile)
	e.GET("/fileServer/getProgress", controler.GetProgress)
	e.Run(":8848")
}
