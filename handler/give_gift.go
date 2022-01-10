package handler

import (
	"encoding/json"
	"net/http"
	"ranking/model"
	"ranking/service"
)

// GiveGift 送礼
func GiveGift(w http.ResponseWriter, req *http.Request) {
	var giftRecDto model.GiftRecDto
	err := json.NewDecoder(req.Body).Decode(&giftRecDto)
	if err != nil {
		Failure(w, err)
		return
	}

	err = service.GiveGift(giftRecDto)
	if err != nil {
		Failure(w, err)
	} else {
		Success(w, nil)
	}
}
