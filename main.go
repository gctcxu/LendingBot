package main

import (
	"lending/kucoin"
	"lending/poloniex"
)

func main() {
	kucoin.KucoinLendingMain()
	poloniex.PoloniexLendingMain()
	select {}
}
