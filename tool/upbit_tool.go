package tool

import (
	"fmt"
	"github.com/jekeun/upbit-go"
	"github.com/jekeun/upbit-go/types"
	"github.com/jekeun/upbit-go/util"
	"log"
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

/*
 * 주문 내역으로부터 해당되는 코인 리스트 구하기
 */
func GetCoinsFromOrders(orders []*types.Order) (coins []string) {
	coins = make([]string, 0)

	for _, value := range orders {
		coins = append(coins, value.Market)
	}

	return
}

/*
 * DayCandle 로부터 현재가 구하기
 * 정보가 없으면 0.0
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
	profit = 0.0

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
			order = value
			bExist = true
			return
		}
	}
	return
}

/*
 * 매도 가능 코인 목록 구하기
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

/*
 * 매도 주문
 */
func AskOrder(client *upbit.Client, coin string, balance string, candle *types.DayCandle, ordType string) {

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
 * 시장가 매도 주문
 */
func AskMarketOrder(client *upbit.Client, coin string, balance string) {

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

/*
 * coin 리스트에 해당하는 Minute Candle 리스트 가져오기
 */
func GetMinuteCandlesByCoins(client *upbit.Client, unit int, coins []string, count int) (
	candleMap map[string][]*types.MinuteCandle) {
	candleMap = make(map[string][]*types.MinuteCandle)
	for _, coin := range coins {
		if coin == "KRW" {
			continue
		}

		candles, err := client.MinuteCandles(unit, coin, map[string]string{
			"count": fmt.Sprintf("%d", count),
		})

		if err != nil {
			log.Println(err)
			continue
		}

		candleMap[coin] = candles
	}

	return
}

/*
 * coin 리스트에 해당하는 현재가 정보 가져오기
 */

func GetTickerByCoins(client *upbit.Client, coins []string) (tickerMap map[string]*types.MinuteCandle) {
	tickerMap = make(map[string]*types.MinuteCandle)
	return
}

/*
 * coin 리스트에 해당하는 Day Candle 리스트 가져오기
 */
func GetDayCandlesByCoins(client *upbit.Client, coins []string, count int) (candleMap map[string][]*types.DayCandle) {
	candleMap = make(map[string][]*types.DayCandle)
	for _, coin := range coins {
		if coin == "KRW" {
			continue
		}

		candles, err := client.DayCandles(coin, map[string]string{
			"count": fmt.Sprintf("%d", count),
		})

		if err != nil {
			log.Println(err)
			continue
		}

		candleMap[coin] = candles
	}

	return
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
