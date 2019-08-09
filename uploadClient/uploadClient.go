package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)
func getFileSize2(filename string) int64 {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}
func uploadFile2(filename string, host string){
	var chunk int64 = 1024*256
	var start int64 = 0
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return
	}

	file_size := fi.Size()
	for file_size > start {
		left := chunk
		if file_size <= start + chunk{
			left = file_size-start
		}
		buf := make([]byte,left)
		fh.ReadAt(buf,start)
		target_url := fmt.Sprintf("%s?file_name=%s&chunk_start=%d&chunk_size=%d",host,filename,start,left)
		res ,err:= postChunk(filename,target_url,bytes.NewBuffer(buf),left)
		if err != nil {
			fmt.Println(err)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(string(body))
		res.Body.Close()
		start = start + chunk
	}
}
func postChunk(filename string, target_url string,file_chunk *bytes.Buffer,buf_size int64) (*http.Response, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}
	boundary := body_writer.Boundary()
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	request_reader := io.MultiReader(body_buf, file_chunk, close_buf)

	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = buf_size + int64(body_buf.Len()) + int64(close_buf.Len())

	return http.DefaultClient.Do(req)
}
func getProgress2(filename string, host string){

	res, err := http.Get(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func main() {
	target_url := "http://192.168.2.56:8899/upload"
	filename := "555.jpg"
	uploadFile2(filename, target_url)

}
