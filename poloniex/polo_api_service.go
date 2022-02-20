package poloniex

import (
	"encoding/json"
	"lending/common"
	"time"

	"strconv"

	"github.com/bitly/go-simplejson"
)

type PoloApiService struct {
	requester *PoloHttpClient
	logger    *common.Logger
}

type LendingOffer struct {
	Rate     float64 `json:"rate,string"`
	Amount   float64 `json:"amount,string"`
	RangeMin int     `json:"rangeMin"`
	RangeMax int     `json:"rangeMax"`
}

type LendingOrder struct {
	Id       int64   `json:"id"`
	Rate     float64 `json:"rate,string"`
	Amount   int     `json:"amount"`
	Duration int     `json:"duration"`
	Currency string  `json:"currency"`
}

func NewPoloApiService(apiKey string, apiSecret string) *PoloApiService {
	logger := common.NewLogger("poloniex.log")
	return &PoloApiService{&PoloHttpClient{apiKey, apiSecret}, logger}
}

func (poloApiService *PoloApiService) GetLendingOfferList(currency string) (lendingOfferList []LendingOffer) {
	res := poloApiService.requester.Get("/public?command=returnLoanOrders&currency=" + currency)

	json_byte, _ := simplejson.NewJson([]byte(res))
	json_array, _ := json_byte.Get("offers").MarshalJSON()
	json.Unmarshal(json_array, &lendingOfferList)

	return lendingOfferList
}

func (poloApiService *PoloApiService) GetLendingBalance(currency string) float64 {
	now := time.Now()

	params := map[string]string{"command": "returnAvailableAccountBalances", "nonce": strconv.FormatInt(now.UnixMilli(), 10)}

	res := poloApiService.requester.Post("/tradingApi", params)

	json, _ := simplejson.NewJson([]byte(res))

	available, _ := json.GetPath("lending", currency).String()

	availableLendingBalance, _ := strconv.ParseFloat(available, 64)
	return availableLendingBalance
}

func (poloApiService *PoloApiService) CreateLendingOrder(currency string, amount int, lendingRate float64, duration int) int64 {
	now := time.Now()

	params := map[string]string{"command": "createLoanOffer", "currency": currency, "amount": strconv.FormatInt(int64(amount), 10), "duration": strconv.Itoa(duration), "lendingRate": strconv.FormatFloat(lendingRate, 'E', 10, 64), "nonce": strconv.FormatInt(now.UnixMilli(), 10)}

	res := poloApiService.requester.Post("/tradingApi", params)

	json, _ := simplejson.NewJson([]byte(res))

	orderId, _ := json.Get("orderID").Int64()

	poloApiService.logger.Log("Order Created:", orderId)

	return orderId
}

func (poloApiService *PoloApiService) CancelLendingOrder(orderNumber int64) {
	now := time.Now()

	params := map[string]string{"command": "cancelLoanOffer", "orderNumber": strconv.FormatInt(orderNumber, 10), "nonce": strconv.FormatInt(now.UnixMilli(), 10)}

	poloApiService.requester.Post("/tradingApi", params)

	poloApiService.logger.Log("Order Canceled:", orderNumber)
}

func (poloApiService *PoloApiService) GetOpenOrder(currency string) (openLendingOrderList []LendingOrder) {

	now := time.Now()

	params := map[string]string{"command": "returnOpenLoanOffers", "nonce": strconv.FormatInt(now.UnixMilli(), 10)}

	res := poloApiService.requester.Post("/tradingApi", params)

	jsonRoot, _ := simplejson.NewJson([]byte(res))
	jsonArray, _ := jsonRoot.Get(currency).Encode()

	json.Unmarshal(jsonArray, &openLendingOrderList)

	return openLendingOrderList
}
