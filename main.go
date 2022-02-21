package main

import (
	"lending/poloniex"
)

func main() {
	//kucoin.KucoinLendingMain()
	poloniex.PoloniexLendingMain()
	select {}
}
