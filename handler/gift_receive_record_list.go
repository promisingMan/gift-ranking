package handler

import (
	"log"
	"net/http"
	"ranking/model"
	"strconv"
)

// GiftRecRecordList 收礼流水查询
func GiftRecRecordList(w http.ResponseWriter, req *http.Request) {
	// 解析入参
	anchorId, page, limit := handleArgs(req)
	result := model.GetGiftRecRecordListByAnchorId(anchorId, page, limit)
	Success(w, result)
}

func handleArgs(req *http.Request) (int, int, int) {
	values := req.URL.Query()
	var anchorIdStr = values.Get("anchorId")
	anchorId, err := strconv.Atoi(anchorIdStr)
	if err != nil {
		log.Panicln("anchorId parse failed", err)
	}
	var pageStr = values.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Panicln("page parse failed", err)
	}
	var limitStr = values.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Panicln("limit parse failed", err)
	}
	return anchorId, page, limit
}
