package poloniex

import (
	"encoding/json"
	"io/ioutil"
	"lending/common"

	"github.com/robfig/cron/v3"
)

type PoloApyScheduler struct {
	LendingOfferList []LendingOffer `json:"lendingOfferList"`
	ApyRecordList    []float64      `json:"apyRecordList"`

	poloApiService *PoloApiService
}

var path = "poloniex.json"

func NewPoloApyScheduler() *PoloApyScheduler {

	scheduler := &PoloApyScheduler{poloApiService: NewPoloApiService(common.AppConfig.Kucoin.ApiKey, common.AppConfig.Poloniex.ApiSecret)}
	scheduler.restore()

	return scheduler
}

func (scheduler *PoloApyScheduler) save() {
	json_byte, _ := json.Marshal(scheduler)
	ioutil.WriteFile(path, json_byte, 0644)
}

func (scheduler *PoloApyScheduler) restore() {
	file_byte, _ := ioutil.ReadFile(path)

	var copy PoloApyScheduler

	json.Unmarshal(file_byte, &copy)

	scheduler.LendingOfferList = copy.LendingOfferList
	scheduler.ApyRecordList = copy.ApyRecordList
}

func (scheduler *PoloApyScheduler) setNewLendingOffer(lendingOfferList []LendingOffer) {

	maxApy := scheduler.calculateDurationMaxApy(lendingOfferList, scheduler.LendingOfferList)

	scheduler.LendingOfferList = lendingOfferList

	if scheduler.ApyRecordList == nil {
		scheduler.ApyRecordList = make([]float64, 0)
	}

	scheduler.ApyRecordList = append(scheduler.ApyRecordList, maxApy)

	if len(scheduler.ApyRecordList) >= 100 {
		scheduler.ApyRecordList = scheduler.ApyRecordList[len(scheduler.ApyRecordList)-100:]
	}

	scheduler.save()
}

func (scheduler *PoloApyScheduler) calculateDurationMaxApy(newLendingOfferList []LendingOffer, oldLendingOfferList []LendingOffer) float64 {

	maxApy := 0.0

	if oldLendingOfferList == nil {
		return maxApy
	}

	var i, j = 0, 0

	for i < len(newLendingOfferList) {
		for j < len(oldLendingOfferList) {
			if newLendingOfferList[i].Rate == oldLendingOfferList[j].Rate { //利率相同, 判斷金額是否有變少, 變少表示有被借出
				if newLendingOfferList[i].Amount < oldLendingOfferList[j].Amount {
					if newLendingOfferList[i].Rate > maxApy {
						maxApy = newLendingOfferList[i].Rate
					}
				}

				i++
				j++
			} else if newLendingOfferList[i].Rate > oldLendingOfferList[j].Rate { //某一項借出 只存在舊的借出列表上, 表示該利率都被借光了
				if oldLendingOfferList[j].Rate > maxApy {
					maxApy = newLendingOfferList[i].Rate
				}

				j++
			} else {
				i++
			}

			if i >= len(newLendingOfferList) || j >= len(oldLendingOfferList) {
				break
			}
		}

		if i >= len(newLendingOfferList) || j >= len(oldLendingOfferList) {
			break
		}
	}

	if maxApy == 0 {
		maxApy = scheduler.ApyRecordList[len(scheduler.ApyRecordList)-1]
	}

	return maxApy
}

func (scheduler *PoloApyScheduler) StartPollingApy() {

	cronRule := "0/10 * * * * *"
	cronScheduler := cron.New(cron.WithSeconds())

	cronScheduler.AddFunc(cronRule, func() {
		lendingOfferList := poloApiService.GetLendingOfferList("USDT")
		scheduler.setNewLendingOffer(lendingOfferList)
	})

	cronScheduler.Start()
}

func (scheduler *PoloApyScheduler) GetMaxPay() float64 {
	maxApy := 0.0

	for i := 0; i < len(scheduler.ApyRecordList); i++ {
		if scheduler.ApyRecordList[i] > maxApy {
			maxApy = scheduler.ApyRecordList[i]
		}
	}

	return maxApy
}

func (scheduler *PoloApyScheduler) GetMinPay() float64 {
	minApy := 1.0

	for i := 0; i < len(scheduler.ApyRecordList); i++ {
		if scheduler.ApyRecordList[i] < minApy {
			minApy = scheduler.ApyRecordList[i]
		}
	}

	return minApy
}
