package main

import (
	"fmt"
	"lending/common"
	"lending/poloniex"
	"testing"
)

func TestPoloGetMarketInfo(t *testing.T) {
	poloHttpClient := poloniex.NewPoloHttpClient(common.AppConfig.Poloniex.ApiKey, common.AppConfig.Poloniex.ApiSecret)

	res := poloHttpClient.Get("/public?command=returnLoanOrders&currency=USDT")

	if len(string(res)) == 0 {
		t.Fail()
	}
}

func TestGetBalance(t *testing.T) {
	poloniexApiService := poloniex.NewPoloApiService(common.AppConfig.Poloniex.ApiKey, common.AppConfig.Poloniex.ApiSecret)

	balance := poloniexApiService.GetLendingBalance("USDT")
	fmt.Println(balance)
}

func TestCreateLoanOffer(t *testing.T) {
	poloniexApiService := poloniex.NewPoloApiService(common.AppConfig.Poloniex.ApiKey, common.AppConfig.Poloniex.ApiSecret)

	poloniexApiService.CreateLendingOrder("USDT", 500, 0.05, 2)
}
