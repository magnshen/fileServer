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
	e.POST("/fileServer/uploadDelete", controler.UploadDelete)
	e.GET("/fileServer/download/:fileName", controler.GetFile)
	e.Run(":8848")
}
