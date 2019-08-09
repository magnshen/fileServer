package model

import (
	"bytes"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)
const TargetUrl = "http://localhost:8899/fileServer"


type UploadModel struct {
	userId string   //用户id  "780002"
	filePath string   //文件路径，包含文件名
	uploadPath string  //上传路径
	uploadName string  //上传文件名
	fileSize int64  //文件总大小
	fileSizeStr string  //文件总大小 字符串类型
	fileHash string  //文件哈希，由上传路径 + 上传文件名 + 大小  计算得到
	progress int64   //进度，已经传了多少
	isReady bool    //如果续传 是否准备好
}



func (self *UploadModel) Init(userId,filePath,uploadPath string)error{
	self.isReady = false
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	self.userId = userId
	self.filePath = filePath
	self.fileSize = fileInfo.Size()
	self.uploadPath = uploadPath
	self.uploadName = path.Base(filePath)
	self.progress = 0

	fileSizeStr := strconv.FormatInt(self.fileSize,10)
	self.fileSizeStr = fileSizeStr
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(fmt.Sprintf("%s-%s-%s",uploadPath,self.uploadName,fileSizeStr)))
	result := Sha1Inst.Sum([]byte(""))
	self.fileHash = base32.StdEncoding.EncodeToString(result)  //上传路径 + 上传文件名 + 大小 计算hash 再使用base32编码转字符串
	return nil
}

func (self *UploadModel) GetProgressFormServer()(int64,error){

	u, _ := url.Parse(TargetUrl+"/getProgress")
	q := u.Query()
	q.Set("user", self.userId)
	q.Set("file_name", self.uploadName)
	q.Set("file_path", self.uploadPath)
	q.Set("file_size", self.fileSizeStr)
	q.Set("file_hash", self.fileHash)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String());

	if err != nil {
		// handle error
		return 0,err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("read response error")
		return 0,err
	}
	resData := &progressResponse{}
	json.Unmarshal(body,resData)
	if resData.Code >= 0 {
		self.progress = resData.Data.Progress
		return self.progress,nil
	}else{
		return 0,errors.New(resData.Description)
	}
}

func (self *UploadModel) UploadStart()error{
	fh, err := os.Open(self.filePath)
	if err != nil {
		fmt.Println("Error opening file")
		return err
	}
	writer := Writer{fh,self.progress}
	u, _ := url.Parse(TargetUrl+"/upload")
	q := u.Query()
	q.Set("user", self.userId)
	q.Set("file_name", self.uploadName)
	q.Set("file_path", self.uploadPath)
	q.Set("file_size", self.fileSizeStr)
	q.Set("file_hash", self.fileHash)
	u.RawQuery = q.Encode()
	apizUrl := u.String()


	header_buf := bytes.NewBufferString("")
	header_buf_writer := multipart.NewWriter(header_buf)

	_, err = header_buf_writer.CreateFormFile("file", self.uploadName)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}
	boundary := header_buf_writer.Boundary()
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	r, w := io.Pipe()
	go writer.doWrite(w)
	request_reader := io.MultiReader(header_buf, r, close_buf)
	req, err := http.NewRequest("POST", apizUrl, request_reader)
	if err != nil {
		fmt.Println("NewRequest Error")
		return err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength =self.fileSize-self.progress + int64(header_buf.Len()) + int64(close_buf.Len())
	resp,err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("request an error")
		return err
	}
	resData := &progressResponse{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(resData)
	if err != nil {
		fmt.Println("json decode error")
		return err
	}
	if resData.Code >= 0 {
		self.progress = resData.Data.Progress
		fmt.Println(resData.Data.Complete)
		fmt.Println("success")
	}else{
		fmt.Println(resData.Description)
	}
	return nil
}