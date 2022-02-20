package main

import (
	"fmt"
	"lending/common"
	"lending/kucoin"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

var config common.Config
var kuApiService *kucoin.KuApiService

func init() {

	currentPath, _ := os.Getwd()

	file, _ := os.ReadFile(currentPath + "/config.yaml")

	yaml.Unmarshal(file, &config)
	kuApiService = kucoin.NewKuApiService(config.Kucoin.ApiKey, config.Kucoin.ApiSecret, config.Kucoin.ApiPassphrase)
}

func TestKuGetBalance(t *testing.T) {
	assetList := kucoin.GetBalance()
	fmt.Println(assetList)
}

func TestKuGet(t *testing.T) {
	kuHttpClient := kucoin.NewKuHttpClient(config.Kucoin.ApiKey, config.Kucoin.ApiSecret, config.Kucoin.ApiPassphrase)

	res := kuHttpClient.Get("/api/v1/accounts")
	fmt.Println(res)
	if len(res) == 0 {
		t.Error("There is no available asset")
	}
}

func TestKuSliceLending(t *testing.T) {
	lendingOrderList := kucoin.SliceLending(24543, 28, 0.05, "USDT")
	fmt.Println(lendingOrderList)
}
