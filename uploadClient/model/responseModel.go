package model

type responseData struct {
	Code int `json:"code"`
	Description string  `json:"description"`
}

type progressData struct {
	FileName string `json:"fileName"`
	Progress int64  `json:"progress"`
	Complete bool `json:"complete"`
}

type progressResponse struct {
	responseData
	Data progressData `json:"data"`
}