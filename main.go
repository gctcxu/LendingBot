package main

import (
	"lending/kucoin"
	"lending/poloniex"
)

func main() {
	kucoin.KucoinLendingMain()
	poloniex.PoloniexLendingMain()
	select {}
	//fmt.Println(5e-05 == 5e-05)
}
