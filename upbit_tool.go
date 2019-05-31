package upbit

import (
	"fmt"
	"github.com/jekeun/upbit-go/types"
	"github.com/jekeun/upbit-go/util"
	"strconv"
	"strings"
)

var DefaultMarket = "KRW"

func GetBalanceMap(balances []*types.Balance) (balanceMap map[string]*types.Balance) {

	balanceMap = make(map[string]*types.Balance)

	for _, value := range balances {

		currency := util.GetMarketFromCurrency(value.Currency, DefaultMarket)

		balanceMap[currency] = value
	}

	return
}

func GetCoinsFromOrders(orders []*types.Order) (coins []string) {
	coins = make([]string, 0)

	for _, value := range orders {
		coins = append(coins, value.Market)
	}

	return
}

/*
 * DayCandle 로부터 currentPrice 구하기
 */
func GetCurrentPriceFromDayCandle(candleInfo []*types.DayCandle) (currentPrice float64){
	currentPrice = 0.0

	if len(candleInfo)  > 0 {
		currentPrice = candleInfo[0].TradePrice
	}

	return
}

/*
 * 수익률 구하기
 */
func GetProfitRate(balance *types.Balance, currentPrice float64) (profit float64) {

	fAvgKrwBuyPrice, err := strconv.ParseFloat(balance.AvgKrwBuyPrice, 64)

	if err != nil {
		return
	}

	profit = (currentPrice/fAvgKrwBuyPrice - 1.0) * 100

	return
}

/*
 * desCoins 리스트에 coin이 있는지 체크
 */
func IsExist(coin string, desCoins []string) (bExist bool) {
	bExist = false

	for _, value := range desCoins {
		if strings.EqualFold(coin, value) {
			bExist = true
			return
		}
	}

	return
}

/*
 * order 목록에 coin 주문이 있는지 확인
 */
func ExistOrder(coin string, orderMap map[string][]*types.Order, orderSide string) (order *types.Order, bExist bool) {
	bExist = false

	orders := orderMap[orderSide]

	for _, value := range orders {
		if strings.EqualFold(value.Market, coin) {
			bExist = true
			order = value
			return
		}
	}

	return
}



/*
 * 매도 가능 잔고
 */
func GetBalanceCoinsCanAsk(balances []*types.Balance) (coins []string) {
	coins = make([]string, 0)

	for _, value := range balances {
		// KRW는 매도 가능 잔고가 아님.
		if value.Currency == "KRW" {
			continue
		}
		currency := util.GetMarketFromCurrency(value.Currency, "KRW")

		s, _ := strconv.ParseFloat(value.Balance, 64)

		if s > 0 {
			coins = append(coins, currency)
		}
	}

	return
}

func AskOrder(client *Client, coin string, balance string, candle *types.DayCandle, ordType string) {

	priceStr :=  fmt.Sprintf("%.8f", candle.TradePrice)
	volumeStr :=  balance

	askOrder := types.OrderInfo{
		Identifier: strconv.Itoa(int(util.TimeStamp())),
		Side:       types.ORDERSIDE_ASK,
		Market:     coin,
		Price:      priceStr,
		Volume:     volumeStr,
		OrdType:    ordType}


	_, err := client.OrderByInfo(askOrder)

	if err != nil {
		fmt.Println("주문 에러")
	}

}

/*
 * 시장가 매도
 */
func AskMarketOrder(client *Client, coin string, balance string) {

	volumeStr :=  balance

	askOrder := types.OrderInfo{
		Identifier: strconv.Itoa(int(util.TimeStamp())),
		Side:       types.ORDERSIDE_ASK,
		Market:     coin,
		Price:      "",
		Volume:     volumeStr,
		OrdType:    types.ORDERTYPE_MARKET}

	_, err := client.OrderByInfo(askOrder)

	if err != nil {
		fmt.Println("주문 에러")
	}

}


func GetPriceCanOrder(price float64) (orderPriceStr string) {

	switch {
	case price >= 2000000:
		orderPriceStr = fmt.Sprintf("%d", int(price/1000) * 1000)
	case price >= 1000000 && price < 2000000:
		orderPriceStr = fmt.Sprintf("%d", int(price/1000) * 1000)
	case price >= 500000 && price < 1000000:
		orderPriceStr = fmt.Sprintf("%d", int(price/100) * 100)
	case price >= 100000 && price < 500000:
		orderPriceStr = fmt.Sprintf("%d", int(price/100) * 100)
	case price >= 10000 && price < 100000:
		orderPriceStr = fmt.Sprintf("%d", int(price/10) * 10)
	case price >= 1000 && price < 10000:
		orderPriceStr = fmt.Sprintf("%d", int(price/10) * 10)
	case price >= 100 && price < 1000:
		orderPriceStr = fmt.Sprintf("%d", int(price))
	case price >= 10 && price < 100:
		orderPriceStr = fmt.Sprintf("%.1f", price)
	case price < 10:
		orderPriceStr = fmt.Sprintf("%.2f", price)
	}

	return
}
