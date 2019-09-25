package arsHash

import (
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
)

const (
	arsMinSizeToMd5 = 100 << 20
)

func bytecopy(dst []byte, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i]
	}
}

func FileHash(path string) (int64, string,error) {
	md5str := ""
	f, err := os.Open(path)
	if err != nil {
		return 0, md5str, err
	}
	defer f.Close()
	size, _ := f.Seek(0, 2)
	f.Seek(0, os.SEEK_SET)
	if size > arsMinSizeToMd5 {
		sizelen := strconv.FormatInt(size, 10)
		sizeNumLen := len(sizelen)
		data := make([]byte, sizeNumLen+512*100+100)
		bytecopy(data, []byte(sizelen))
		avg := size / 512
		//buf := make([]byte, 100)
		for i := 0; i < 512; i++ {
			//off := (int64)(i) * avg
			f.ReadAt(data[sizeNumLen+i*100:sizeNumLen+i*100+100], (int64)(i)*avg)
			//bytecopy(data[len(sizelen)+int(i*100):], buf)
		}
		f.Seek(size-100-1, os.SEEK_SET)
		num ,_:= f.Read(data[sizeNumLen+512*100:])
		fmt.Println(num)
		has := md5.Sum(data)
		md5str = fmt.Sprintf("%x", has)

	} else {
		//fmt.Println("Size:", size, " < ", ArsMinSizeToMd5)
		data := make([]byte, size)
		f.Read(data)
		has := md5.Sum(data)
		md5str = fmt.Sprintf("%x", has)
	}
	//fmt.Println("md5str:", md5str)
	return size,md5str,nil
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
