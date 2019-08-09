package model

import (
	"io"
	"os"
)

type Writer struct {
	fileHandler *os.File   //这是我要读取的句柄,本地文件
	startPoin int64    //断点续传位置
}
//这是做实验
func (w *Writer)doWrite(writer *io.PipeWriter){
	buf:=make([]byte,1024)       //每次读取大小，设置太小会影响速度。太大也没用，瓶颈是网络，而且增加内存
	w.fileHandler.Seek(w.startPoin,1)
	m := int64(0)
	for {
		n,err := w.fileHandler.Read(buf)
		if n>0{
			writer.Write(buf[0:n])
		}
		if m > 48<<20{
			//time.Sleep(16*time.Second)
			writer.Write(nil)
			break
		}
		m = m+ int64(n)
		if err==io.EOF {//结束
			break
		}
	}
	defer writer.Close()
}

//func (w *Writer)doWrite(writer *io.PipeWriter){
//	buf:=make([]byte,1024)       //每次读取大小，设置太小会影响速度。太大也没用，瓶颈是网络，而且增加内存
//	w.fileHandler.Seek(w.startPoin,1)
//	for {
//		n,err := w.fileHandler.Read(buf)
//		if n>0{
//			writer.Write(buf[0:n])
//		}
//		if err==io.EOF {//结束
//			fmt.Println()
//			break
//		}
//	}
//	defer writer.Close()
//}