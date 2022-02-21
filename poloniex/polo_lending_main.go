package poloniex

import (
	"lending/common"
	"math"
	"os"

	"github.com/robfig/cron/v3"
)

var poloApiService *PoloApiService
var poloniexApyScheduler *PoloApyScheduler
var isProduction bool = os.Getenv("OS") != "Windows_NT"

func init() {

	poloApiService = NewPoloApiService(common.AppConfig.Poloniex.ApiKey, common.AppConfig.Poloniex.ApiSecret)
	poloniexApyScheduler = NewPoloApyScheduler()
}

func SliceLending(balance float64, targetTerm int, targetDailyIntRate float64, currency string) []LendingOrder {

	remaining := int(balance)

	toSubmitLendingOrder := make([]LendingOrder, 0)

	for remaining >= common.AppConfig.Kucoin.SliceSize {
		var sliceSize int

		if remaining < 2*common.AppConfig.Poloniex.SliceSize {
			sliceSize = remaining
		} else {
			sliceSize = common.AppConfig.Poloniex.SliceSize
		}

		remaining -= sliceSize

		toSubmitLendingOrder = append(toSubmitLendingOrder, LendingOrder{Duration: targetTerm, Amount: sliceSize, Rate: targetDailyIntRate, Currency: currency})
	}

	if remaining >= 1 {
		toSubmitLendingOrder = append(toSubmitLendingOrder, LendingOrder{Duration: targetTerm, Amount: remaining, Rate: targetDailyIntRate, Currency: currency})
	}

	return toSubmitLendingOrder
}

func BatchCancelLendingOrder(lendingOrderList []LendingOrder) {
	for _, openOrder := range lendingOrderList {
		poloApiService.CancelLendingOrder(openOrder.Id)
	}
}

func BatchSubmitLendingOrder(lendingOrderList []LendingOrder) {
	for _, lendingOrder := range lendingOrderList {
		poloApiService.CreateLendingOrder(lendingOrder.Currency, lendingOrder.Amount, lendingOrder.Rate, lendingOrder.Duration)
	}
}

func PoloniexLendingMain() {

	var cronRule string
	if isProduction {
		cronRule = "0 0/10 * * * *"
	} else {
		cronRule = "0/10 * * * * *"
	}

	poloniexApyScheduler.StartPollingApy()

	scheduler := cron.New(cron.WithSeconds())
	scheduler.AddFunc(cronRule, func() {
		balance := poloApiService.GetLendingBalance(common.AppConfig.Poloniex.LendingCurrency)
		if balance < 1 {
			poloApiService.logger.Log("Balance is not enough, check next time")
			return
		}

		targetDailyRate := poloniexApyScheduler.GetAvgApy()
		if targetDailyRate == 0 {
			poloApiService.logger.Log("Daily Rate 0, there must be something wrong")
			return
		}

		//取得借出掛單
		openLendingOrderList := poloApiService.GetOpenOrder(common.AppConfig.Poloniex.LendingCurrency)

		//最少分段數量
		minNumPiece := math.Ceil(balance / float64(common.AppConfig.Kucoin.SliceSize))

		if len(openLendingOrderList) > 0 {
			//掛單數量為最少分段數量 且 最佳利率沒改變 => 則不進行動作
			if int(minNumPiece) == len(openLendingOrderList) && openLendingOrderList[0].Rate == targetDailyRate {
				poloApiService.logger.Log("There is no better interest, wait for next time")
				return
			}

			//無任何的區塊可以借出, 直接中止
		} else if minNumPiece == 0 {
			poloApiService.logger.Log("There is lending slice, wait for next time")
			return
		}

		BatchCancelLendingOrder(openLendingOrderList)
		balance = poloApiService.GetLendingBalance(common.AppConfig.Poloniex.LendingCurrency)
		toSubmitedLendingOrder := SliceLending(balance, 2, targetDailyRate, common.AppConfig.Poloniex.LendingCurrency)

		BatchSubmitLendingOrder(toSubmitedLendingOrder)
	})

	scheduler.Start()
}
