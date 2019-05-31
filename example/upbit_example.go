package main

import (
	"fmt"
	"github.com/jekeun/upbit-go"
	"github.com/jekeun/upbit-go/util"
	"log"
	"strconv"
	"time"
)

var accessKey = "Your Access Key"
var secretKey = "Your Secret Key"

func main() {
//	getMarkets()
	getAccounts()


	for {
		// getMinutesCandles()
		getDayCandles()
		// getTicker()
		// getOrderChance()
		//fmt.Println("주문")
		//postOrder()
		//
		//time.Sleep(time.Duration(time.Second * 5))
		//
		//fmt.Println("주문 취소")
		getOrders()


		// orderCancel()
		time.Sleep(time.Duration(time.Second * 5))
	}

	fmt.Scanln()
}

func getOrders() {
	client := upbit.NewClient(accessKey, secretKey)

	orders, err := client.Orders("KRW-ATOM","wait", 1, "desc")

	if err != nil {
		log.Panicln(err)
	}

	if len(orders) > 0 {
		fmt.Println(orders[0])

		order, err := client.CancelOrder(orders[0].Uuid)

		if err != nil {
			log.Panicln(err)
		}

		fmt.Println(order)
	}

}

func orderCancel() {

}


func postOrder() {
	client := upbit.NewClient(accessKey, secretKey)

	order, err := client.Order(
		strconv.Itoa(int(util.TimeStamp())),
		"ask",
		"KRW-ATOM",
		"6700",
		"10",
		"limit",
	)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println(order)


}

func getTicker() {
	client := upbit.NewClient(accessKey, secretKey)

	ticker, err := client.Ticker("KRW-XRP")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(ticker[0])

}

func getDayCandles() {
	client := upbit.NewClient(accessKey, secretKey)

	candles, err := client.DayCandles("KRW-XRP", map[string]string{
		"count": "60",
	})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("===")
	fmt.Printf("Opening %f \n", candles[0].OpeningPrice)
	fmt.Printf("High %f \n", candles[0].HighPrice)
	fmt.Printf("Low %f \n", candles[0].LowPrice)
	fmt.Printf("Trade %f \n", candles[0].TradePrice)

	//fmt.Printf("Opening %f \n", candles[1].OpeningPrice)
	//fmt.Printf("High %f \n", candles[1].HighPrice)
	//fmt.Printf("Low %f \n", candles[1].LowPrice)
	//fmt.Printf("Trade %f \n", candles[1].TradePrice)




	// fmt.Println(candles[1])

}

func getOrderChance() {
	client := upbit.NewClient(accessKey, secretKey)
	orderChance, err := client.OrderChance("KRW-ATOM")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(orderChance)
}


func getMinutesCandles() {
	client := upbit.NewClient(accessKey, secretKey)

	candles, err := client.MinuteCandles(60, "KRW-XRP", map[string]string{
		"count": "60",
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(candles))
	fmt.Println(candles[0])

	//for _, value := range candles {
	//	fmt.Println(value)
	//}



}

func getMarkets() {
	client := upbit.NewClient(accessKey,secretKey)

	markets, err := client.Markets()
	if err != nil {
		return
	}

	fmt.Println(markets[0].Market)
}

func getAccounts() {
	client := upbit.NewClient(accessKey,secretKey)

	accounts, err := client.Accounts()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(accounts))

	for _, value := range accounts {
		fmt.Println( value.Currency + " : " + value.Balance)
	}
	// Output:
	// KRW
}
