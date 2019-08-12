package controler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)
const UserDataPath = "/data/cloud/data/data"
const TempDataPath = "/data/cloud/data/data/temp"
func getFileSize(filename string) int64 {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}
func pathExists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func createFilePath(path string)error{
	if !pathExists(path) {
		err := os.MkdirAll(path,os.ModePerm)
		return err
	}
	return nil

}

//func filePathToHash(path string) string{
//	Sha1Inst:=sha1.New()
//	Sha1Inst.Write([]byte(path))
//	result := Sha1Inst.Sum([]byte(""))
//	return base32.StdEncoding.EncodeToString(result)
//}

//func creatUploadCnf(cnfFile string ,fileSize int64) error{
//	f, err := os.OpenFile(cnfFile, os.O_WRONLY|os.O_TRUNC, 0666)  //写入时会覆盖文件
//	defer f.Close()
//	if err != nil{
//		return err
//	}
//	string := strconv.FormatInt(fileSize,10)
//	_, err = f.WriteString(string)
//	return err
//}

func getFileNameFormRepeatNane(filePath,fileName string)(string ,error){
	i := 0
	var file string
	fileSuffix := path.Ext(fileName)
	filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
	for{
		if i ==0{
			file = fileName
		}else{
			file = fmt.Sprintf("%s(%d)%s",filenameOnly,i,fileSuffix)
		}
		isExis := pathExists(filePath+file)
		if !isExis{
			break
		}
		i++
	}
	return file ,nil
}

//func getFileSizeFromCnf(cnfFile string) (int64 ,error){
//	fp, err := os.Open(cnfFile)
//	if err != nil {
//		return 0, err
//	}
//	defer fp.Close()
//	buffer := make([]byte, 16)
//
//	_, err = fp.Read(buffer)
//	if err != nil {
//		return 0, err
//	}
//	fileSize, err := strconv.ParseInt(string(buffer), 10, 64)  //读出来转字符串再转int64
//	if err != nil {
//		return 0, err
//	}
//	return fileSize,nil
//}

func GetProgress(c *gin.Context) {
	user := c.Query("user")
	fileName := c.Query("file_name")
	filePath := c.Query("file_path")
	fileHash := c.Query("file_hash")  //可以用uuid
	fileTmp := UserDataPath+"/tmp/"+user+"/"+fileHash
	if fileHash == ""{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"file_hash is null"})
		return
	}
	var progress int64 = 0
	isExists := pathExists(fileTmp)
	if isExists{
		progress = getFileSize(fileTmp)
	}
	newFile,err := getFileNameFormRepeatNane(UserDataPath+"/User/"+user+filePath+"/",fileName)
	if err != nil{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"get new file name error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0,"description":"","data":gin.H{"progress": progress,"fileName":newFile}})
}



func AppendHandle(c *gin.Context) {
	user := c.Query("user")
	fileName := c.Query("file_name")
	filePath := c.Query("file_path")
	fileHash := c.Query("file_hash")  //可以用uuid
	fileSizeStr := c.Query("file_size")
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil ||fileSize <0{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"file_size must be uint64"})
		return
	}
	if fileHash == ""{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"file_hash is null"})
		return
	}
	fileTempPath := TempDataPath+"/"+user+"/"+fileHash
	err = createFilePath(TempDataPath+"/"+user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"creat tmp folder failed"})
		return
	}
	fileTemp, err := os.OpenFile(fileTempPath, os.O_CREATE|os.O_RDWR, 0666)
	defer fileTemp.Close()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"creat tmp file failed"})
		return
	}
	curSize, _ := fileTemp.Seek(0, 2)  //文件句柄跳到最后，并返回偏移量(临时文件大小)
	//fileUpload, err := c.FormFile("file")
	//fileUpload, header, err := c.Request.FormFile("file")   //读取header,可以从range 鉴定临时文件大小是否一致
	fileUpload := c.Request.Body
	defer fileUpload.Close()

	buf := make([]byte,1024)
	for {
		n,err := fileUpload.Read(buf)   //网络原因,每次读不一定是1024
		if n>0{
			fileTemp.Write(buf[0:n])
		}
		curSize = curSize + int64(n)
		if err==io.EOF { //结束
			fileTemp.Close()
			fmt.Println("finish ?")
			break
		}
	}
	if curSize < fileSize{
		c.JSON(http.StatusOK, gin.H{"code": 0,"description":"incomplete","data":gin.H{"fileName":fileName,"progress": curSize,"complete":false}})
		return
	}
	if curSize > fileSize{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"Temp file is bigger then file size"})
		return
	}

	physicsPath := fmt.Sprintf("%s/User/%s%s/",UserDataPath,user,filePath)
	err = createFilePath(physicsPath)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"Create User Upload Path failed"})
		return
	}
	newFileName,err := getFileNameFormRepeatNane(physicsPath,fileName)
	os.Rename(fileTempPath, physicsPath+newFileName)
	c.JSON(http.StatusOK, gin.H{"code": 0,"description":"success","data":gin.H{"fileName":newFileName,"progress": curSize,"complete":true}})
}

func UploadNewFile(c *gin.Context) {
	user := c.Query("user")
	fileName := c.Query("file_name")
	filePath := c.Query("file_path")
	fileHash := c.Query("file_hash")  //可以用uuid
	fileSizeStr := c.Query("file_size")
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil ||fileSize <0{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"file_size must be uint64"})
		return
	}
	if fileHash == ""{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"file_hash is null"})
		return
	}
	fileTempPath := TempDataPath+"/"+user+"/"+fileHash
	err = createFilePath(TempDataPath+"/"+user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"creat tmp folder failed"})
		return
	}
	fileTemp, err := os.OpenFile(fileTempPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer fileTemp.Close()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"creat tmp file failed"})
		return
	}
	curSize := int64(0)
	fileUpload := c.Request.Body

	defer fileUpload.Close()

	buf := make([]byte,1024)
	for {
		n,err := fileUpload.Read(buf)   //网络原因,每次读不一定是1024
		if n>0{
			fileTemp.Write(buf[0:n])
		}
		curSize = curSize + int64(n)
		if err==io.EOF { //结束
			fileTemp.Close()
			fmt.Println("finish ?")
			break
		}
	}
	if curSize < fileSize{
		c.JSON(http.StatusOK, gin.H{"code": 0,"description":"incomplete","data":gin.H{"fileName":fileName,"progress": curSize,"complete":false}})
		return
	}
	if curSize > fileSize{
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"Temp file is bigger then file size"})
		return
	}

	physicsPath := fmt.Sprintf("%s/User/%s%s/",UserDataPath,user,filePath)
	err = createFilePath(physicsPath)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1,"description":"Create User Upload Path failed"})
		return
	}
	newFileName,err := getFileNameFormRepeatNane(physicsPath,fileName)
	os.Rename(fileTempPath, physicsPath+newFileName)
	c.JSON(http.StatusOK, gin.H{"code": 0,"description":"success","data":gin.H{"fileName":newFileName,"progress": curSize,"complete":true}})
}

func init(){
	exist := pathExists(TempDataPath)
	if !exist  {
		// 创建文件夹
		err := os.Mkdir(TempDataPath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			os.Exit(2)
		}
	}
}
