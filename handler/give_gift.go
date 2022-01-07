package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"ranking/model"
	"ranking/service"
)

// GiveGift 送礼
func GiveGift(w http.ResponseWriter, req *http.Request) {
	// 解析json入参
	decoder := json.NewDecoder(req.Body)
	var giftRecDto model.GiftRecDto
	err := decoder.Decode(&giftRecDto)
	if err != nil {
		log.Panicln("parse json input parameters failed", err)
	}
	service.GiveGift(giftRecDto)
	Success(w, nil)
}
