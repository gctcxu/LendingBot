package kucoin

import (
	"math"
	"os"
	"time"

	"lending/common"

	"github.com/robfig/cron/v3"
)

var kuApiService *KuApiService
var isProduction bool = os.Getenv("OS") != "Windows_NT"

type MarketInfo struct {
	LongTerm  LendingOrder
	MidTerm   LendingOrder
	ShortTerm LendingOrder
}

func init() {
	kuApiService = NewKuApiService(common.AppConfig.Kucoin.ApiKey, common.AppConfig.Kucoin.ApiSecret, common.AppConfig.Kucoin.ApiPassphrase)
}

func GetBalance() float64 {

	assetList := kuApiService.GetBalance()

	for _, asset := range assetList {
		if asset.Type == "main" {
			return asset.Available
		}
	}

	return 0
}

func SliceLending(balance float64, targetTerm int, targetDailyIntRate float64, currency string) []LendingOrder {

	remaining := int(balance)

	toSubmitLendingOrder := make([]LendingOrder, 0)

	for remaining >= common.AppConfig.Kucoin.SliceSize {
		var sliceSize int

		if remaining < 2*common.AppConfig.Kucoin.SliceSize {
			sliceSize = remaining
		} else {
			sliceSize = common.AppConfig.Kucoin.SliceSize
		}

		remaining -= sliceSize
		dailyIntRate := targetDailyIntRate

		toSubmitLendingOrder = append(toSubmitLendingOrder, LendingOrder{Term: targetTerm, Size: sliceSize, DailyIntRate: dailyIntRate, Currency: currency})
	}

	if remaining >= 1 {
		toSubmitLendingOrder = append(toSubmitLendingOrder, LendingOrder{Term: targetTerm, Size: remaining, DailyIntRate: targetDailyIntRate, Currency: currency})
	}

	return toSubmitLendingOrder
}

/*找出最近N時間區間內, 最好的利率*/
func GetBestMarketInfo(currency string) MarketInfo {

	lendingOrderList := kuApiService.GetMarketFilledLendingOrder(currency)

	marketInfo := MarketInfo{}

	for _, lendingOrder := range lendingOrderList {
		if time.Now().UnixNano()-lendingOrder.Timestamp > 60*60*1000000000 {
			continue
		}

		if lendingOrder.Term == 7 {
			if marketInfo.ShortTerm == (LendingOrder{}) {
				marketInfo.ShortTerm = lendingOrder
			} else if lendingOrder.DailyIntRate > marketInfo.ShortTerm.DailyIntRate {
				marketInfo.ShortTerm = lendingOrder
			}
		} else if lendingOrder.Term == 14 {
			if marketInfo.MidTerm == (LendingOrder{}) {
				marketInfo.MidTerm = lendingOrder
			} else if lendingOrder.DailyIntRate > marketInfo.MidTerm.DailyIntRate {
				marketInfo.MidTerm = lendingOrder
			}
		} else {
			if marketInfo.LongTerm == (LendingOrder{}) {
				marketInfo.LongTerm = lendingOrder
			} else if lendingOrder.DailyIntRate > marketInfo.LongTerm.DailyIntRate {
				marketInfo.LongTerm = lendingOrder
			}
		}
	}

	return marketInfo
}

func BatchCancelLendingOrder(activeLendingOrderList []ActiveLendingOrder) {
	for _, activeLendingOrder := range activeLendingOrderList {
		kuApiService.CancelLendingOrder(activeLendingOrder.OrderId)
	}
}

func BatchSubmitLendingOrder(lendingOrderList []LendingOrder) {
	for _, lendingOrder := range lendingOrderList {
		kuApiService.CreateLendingOrder(lendingOrder.Currency, lendingOrder.Size, lendingOrder.DailyIntRate, lendingOrder.Term)
	}
}

func KucoinLendingMain() {

	var cronRule string
	if isProduction {
		cronRule = "0 0/10 * * * *"
	} else {
		cronRule = "0/10 * * * * *"
	}

	scheduler := cron.New(cron.WithSeconds())

	scheduler.AddFunc(cronRule, func() {
		balance := GetBalance()
		if balance == 0 {
			kuApiService.logger.Log("Balance 0, Check next time")
			return
		}

		marketInfo := GetBestMarketInfo(common.AppConfig.Kucoin.LendingCurrency)
		targetDailyIntRate := marketInfo.ShortTerm.DailyIntRate

		if targetDailyIntRate == 0 {
			kuApiService.logger.Log("Daily Rate 0, there must be something wrong")
			return
		}

		//取得借出掛單
		activeLendingOrderList := kuApiService.GetActiveLendingOrder(common.AppConfig.Kucoin.LendingCurrency)

		//最少分段數量
		minNumPiece := math.Ceil(balance / float64(common.AppConfig.Kucoin.SliceSize))

		if len(activeLendingOrderList) > 0 {
			//掛單數量為最少分段數量 且 最佳利率沒改變 => 則不進行動作
			if int(minNumPiece) == len(activeLendingOrderList) && activeLendingOrderList[0].DailyIntRate == targetDailyIntRate {
				kuApiService.logger.Log("There is no better interest, wait for next time")
				return
			}

			//無任何的區塊可以借出, 直接中止
		} else if minNumPiece == 0 {
			kuApiService.logger.Log("There is lending slice, wait for next time")
			return
		}

		BatchCancelLendingOrder(activeLendingOrderList)
		balance = GetBalance()
		openLendingOrderList := SliceLending(balance, common.AppConfig.Kucoin.Term, targetDailyIntRate, common.AppConfig.Kucoin.LendingCurrency)

		BatchSubmitLendingOrder(openLendingOrderList)
	})

	scheduler.Start()
	select {}
}
