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
	NewName string `json:"newName"`
	FileSize int64 `json:"fileSize"`
	FileHash string `json:"fileHash"`
}
type progressData struct {
	Progress int64 `json:"progress"`
	FileInfo fileInfo `json:"fileInfo"`
}
type progressResponse struct {
	Code int `json:"code"`
	Description string  `json:"description"`
	Data progressData `json:"data"`
}