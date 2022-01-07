package router

import (
	"net/http"
	"ranking/handler"
)

func RegisterRoutes() {
	http.HandleFunc("/giveGift", handler.GiveGift)
	http.HandleFunc("/rank", handler.Rank)
	http.HandleFunc("/giftRecRecordList", handler.GiftRecRecordList)
}
