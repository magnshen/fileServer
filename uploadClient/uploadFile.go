package main

import (
	"fmt"
	"./model"
)

const localPath = "/Users/you/Documents/GitHub/fileServer/uploadClient"


func main() {
	user := "780001"
	filePath := localPath+"/123.txt"
	uploadPath := "/home"

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
	fmt.Printf("获取上传进度: %d\n",tempSize)
	err = uploadModel.UploadStart()
	//err = uploadModel.UploadDelete()
	if err != nil{
		fmt.Println(err)
		fmt.Println("upload failed")
		return
	}
}
