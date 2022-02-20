package kucoin

import (
	"encoding/json"
	"fmt"
	"lending/common"

	"github.com/bitly/go-simplejson"
)

type Asset struct {
	Available float64 `json:"available,string"`
	Balance   float64 `json:"balance,string"`
	Currency  string  `json:"currency"`
	Holds     string  `json:"holds"`
	Id        string  `json:"id"`
	Type      string  `json:"type"`
}

type LendingOrder struct {
	Currency     string  `json:"currency"`
	Size         int     `json:"size"`
	DailyIntRate float64 `json:"dailyIntRate,string"`
	Term         int     `json:"term"`
	Timestamp    int64   `json:"timestamp"`
	TradeId      string  `json:"tradeId"`
}

type ActiveLendingOrder struct {
	OrderId      string  `json:"orderId"`
	Currency     string  `json:"currency"`
	Size         int     `json:"size,string"`
	FilledSize   int     `json:"filledSize,string"`
	DailyIntRate float64 `json:"dailyIntRate,string"`
	Term         int     `json:"term"`
	CreatedAt    int64   `json:"createdAt"`
}

type KuApiService struct {
	requester *KuHttpClient
	logger    *common.Logger
}

func NewKuApiService(apiKey string, apiSecret string, apiPassphrase string) *KuApiService {
	return &KuApiService{NewKuHttpClient(apiKey, apiSecret, apiPassphrase), common.NewLogger("kucoin.log")}
}

func (kuApiService *KuApiService) GetBalance() []Asset {
	res := kuApiService.requester.Get("/api/v1/accounts")

	var assetList []Asset

	json.Unmarshal([]byte(res), &assetList)

	return assetList
}

func (kuApiService *KuApiService) GetMarketFilledLendingOrder(currency string) []LendingOrder {
	res := kuApiService.requester.Get("/api/v1/margin/trade/last?currency=" + currency)

	var lendingOrderList []LendingOrder

	json.Unmarshal([]byte(res), &lendingOrderList)

	return lendingOrderList
}

func (kuApiService *KuApiService) GetActiveLendingOrder(currency string) (activeLendingOrderList []ActiveLendingOrder) {
	res := kuApiService.requester.Get("/api/v1/margin/lend/active?currency=" + currency + "&currentPage=1&pageSize=50")

	json_root, _ := simplejson.NewJson([]byte(res))
	json_byte_array, _ := json_root.Get("items").Encode()
	json.Unmarshal(json_byte_array, &activeLendingOrderList)
	return
}

func (kuApiService *KuApiService) CancelLendingOrder(orderId string) {
	kuApiService.requester.Delete("/api/v1/margin/lend/" + orderId)
	kuApiService.logger.Log("Order Canceled:", orderId)
}

func (kuApiService *KuApiService) CreateLendingOrder(currency string, size int, dailyIntRate float64, term int) {
	params := map[string]interface{}{
		"currency":     currency,
		"size":         size,
		"dailyIntRate": dailyIntRate,
		"term":         term,
	}

	res := kuApiService.requester.Post("/api/v1/margin/lend", params)
	fmt.Println(res)
	kuApiService.logger.Log("Order Created:", params)
}
