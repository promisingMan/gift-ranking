package handler

import (
	"encoding/json"
	"net/http"
)

// Response 通用响应体，后续可以提取code和msg错误表
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Success 统一格式返回成功response body
func Success(w http.ResponseWriter, data interface{}) {
	response := Response{Code: 0, Data: data}
	_ = json.NewEncoder(w).Encode(response)
}

// Failure 统一格式返回失败response body
func Failure(w http.ResponseWriter, data interface{}) {
	response := Response{Code: 500, Data: data}
	_ = json.NewEncoder(w).Encode(response)
}
