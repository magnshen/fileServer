package main

import (
	"./model"
	"fmt"
)

const localPath = "/Users/you/Documents/GitHub/fileServer/uploadClient"


func main() {
	user := "780001"
	filePath := localPath+"/SNH48-梦想岛.mp4"
	//filePath := localPath+"/123.txt"
	uploadPath := "/*home*/下载"
	uploadModel := model.UploadModel{}
	err := uploadModel.Init(user,filePath,uploadPath)
	if err != nil{
		fmt.Println(err)
		return
	}
	uploadModel.IsCover = true  //是否覆盖上传
	//fileHash := arsHash.FileHash(filePath)
	//fmt.Println(fileHash)

	progressInfo,err := uploadModel.GetProgressFromServer()
	if err != nil{
		fmt.Println(err)
		fmt.Println("get Progress failed")
		return
	}
	fmt.Printf("获取上传进度(缓存文件大小): %d\n",progressInfo.Progress)
	if progressInfo.FileInfo.FileSize > 0{
		fmt.Printf("目标路径有重名文件\n")
		fmt.Printf("远端文件大小: %d\n",progressInfo.FileInfo.FileSize)
		fmt.Printf("远端文件哈希: %s\n",progressInfo.FileInfo.FileHash)
		fmt.Printf("新文件名称: %s\n",progressInfo.FileInfo.NewName)
	}
	err = uploadModel.UploadStart()
	//err = uploadModel.UploadDelete()
	if err != nil{
		fmt.Println(err)
		fmt.Println("upload failed")
		return
	}
}
