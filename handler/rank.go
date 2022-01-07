package handler

import (
	"net/http"
	"ranking/service"
)

// Rank 获取排行榜
func Rank(w http.ResponseWriter, req *http.Request) {
	resp, err := service.GetRank()
	if err != nil {
		Failure(w, err)
	} else {
		Success(w, resp)
	}
}
