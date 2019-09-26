package model

type uploadResponseData struct {
	FileName string `json:"fileName"`
	Progress int64  `json:"progress"`
	Complete bool `json:"complete"`
}

type uploadResponse struct {
	Code int `json:"code"`
	Description string  `json:"description"`
	Data uploadResponseData `json:"data"`
}

type fileInfo struct {
	FileName string `json:"fileName"`
	FileSize int64 `json:"fileSize"`
	FileHash string `json:"fileHash"`
}
type progressData struct {
	NewName string `json:"newName"`
	Progress int64 `json:"progress"`
	FileInfoList []fileInfo `json:"fileInfoList"`
}
type progressResponse struct {
	Code int `json:"code"`
	Description string  `json:"description"`
	Data progressData `json:"data"`
}