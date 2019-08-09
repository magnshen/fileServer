package main

import (
	"fmt"
	"uploadClient/model"
)

const localPath = "C:/GoWork/src/uploadClient"


func main() {
	user := "780001"
	filePath := localPath+"/FF.mp4"
	uploadPath := "/home/xiazai"

	uploadModel := model.UploadModel{}
	err := uploadModel.Init(user,filePath,uploadPath)
	if err != nil{
		fmt.Println(err)
		return
	}

	tempSize,err := uploadModel.GetProgressFormServer()

	if err != nil{
		fmt.Println(err)
		fmt.Println("get Progress failed")
		return
	}
	fmt.Println(tempSize)
	err = uploadModel.UploadStart()
	if err != nil{
		fmt.Println(err)
		fmt.Println("upload failed")
		return
	}
}
