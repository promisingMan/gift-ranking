package handler

import (
	"log"
	"net/http"
	"ranking/model"
	"strconv"
)

// GiftRecRecordList 收礼流水查询
func GiftRecRecordList(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			Failure(w, err)
		}
	}()
	// 解析并校验入参
	anchorId, page, limit := parseAndVerifyArgs(req)
	result, err := model.GetGiftRecRecordListByAnchorId(anchorId, page, limit)
	if err != nil {
		Failure(w, err)
	} else {
		Success(w, result)
	}
}

func parseAndVerifyArgs(req *http.Request) (anchorId, page, limit int) {
	defaultPage := 1
	defaultLimit := 10

	values := req.URL.Query()

	var anchorIdStr = values.Get("anchorId")
	anchorId, err := strconv.Atoi(anchorIdStr)
	if err != nil {
		log.Panic("wrong parameter anchorId")
	}

	var pageStr = values.Get("page")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		log.Panic("wrong parameter page")
	}
	if page <= 0 {
		page = defaultPage
	}

	var limitStr = values.Get("limit")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		log.Panic("wrong parameter limit")
	}
	if limit <= 0 {
		limit = defaultLimit
	}

	return anchorId, page, limit
}
