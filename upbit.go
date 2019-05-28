package upbit

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jekeun/upbit-go/types"
	"github.com/jekeun/upbit-go/util"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseUrl = "https://api.upbit.com/v1"
)

type InvalidParams struct {
	message string
	Err     error
}

func (e *InvalidParams) Error() string {
	return e.message
}

type Client struct {
	accessKey string
	secretKey string
}

func (client *Client) Token(query map[string]string) (tokenStr string, err error) {
	claim := jwt.MapClaims{
		"access_key": client.accessKey,
		"nonce":      util.TimeStamp(),
	}

	if query != nil {
		url := new(url.URL)

		q := url.Query()

		for i, value := range query {
			q.Add(i, value)
		}

		rawQuery := q.Encode()

		claim["query"] = rawQuery
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenStr, err = token.SignedString([]byte(client.secretKey[:]))
	if err != nil {
		return
	}

	return
}

func (client *Client) Markets() (markets []*types.Market, err error) {
	options := &util.RequestOptions{Url: baseUrl + "/market/all"}
	err = util.Request(options, &markets)
	return
}

func (client *Client) MinuteCandles(
	unit int,
	market string,
	params ...map[string]string,
) (candles []*types.MinuteCandle, err error) {
	if !isValidMinuteCandleUnit(unit) {
		err = &InvalidParams{
			message: "Invalid unit",
		}
		return
	}

	query := map[string]string{
		"market": market,
	}

	if len(params) > 0 {
		for _, param := range params {
			for index, value := range param {
				query[index] = value
			}
		}
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/candles/minutes/" + strconv.Itoa(unit),
		Query: query,
	}
	err = util.Request(options, &candles)
	return
}

func (client *Client) DayCandles(
	market string,
	params ...map[string]string,
) (candles []*types.DayCandle, err error) {
	query := map[string]string{
		"market": market,
	}

	if len(params) > 0 {
		for _, param := range params {
			for index, value := range param {
				query[index] = value
			}
		}
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/candles/days",
		Query: query,
	}
	err = util.Request(options, &candles)
	return
}

func (client *Client) WeekCandles(
	market string,
	params ...map[string]string,
) (candles []*types.WeekCandle, err error) {
	query := map[string]string{
		"market": market,
	}

	if len(params) > 0 {
		for _, param := range params {
			for index, value := range param {
				query[index] = value
			}
		}
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/candles/weeks",
		Query: query,
	}
	err = util.Request(options, &candles)
	return
}

func (client *Client) MonthCandles(
	market string,
	params ...map[string]string,
) (candles []*types.MonthCandle, err error) {
	query := map[string]string{
		"market": market,
	}

	if len(params) > 0 {
		for _, param := range params {
			for index, value := range param {
				query[index] = value
			}
		}
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/candles/months",
		Query: query,
	}
	err = util.Request(options, &candles)
	return
}

func (client *Client) TradeTicks(
	market string,
	params ...map[string]string,
) (tradeTicks []*types.TradeTicks, err error) {
	query := map[string]string{
		"market": market,
	}

	if len(params) > 0 {
		for _, param := range params {
			for index, value := range param {
				query[index] = value
			}
		}
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/trades/ticks",
		Query: query,
	}
	err = util.Request(options, &tradeTicks)
	return
}

func (client *Client) Ticker(
	markets string,
) (tickers []*types.Ticker, err error) {
	query := map[string]string{
		"markets": markets,
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/ticker",
		Query: query,
	}
	err = util.Request(options, &tickers)
	return
}

func (client *Client) Orderbooks(
	markets string,
) (orderbooks []*types.Orderbook, err error) {
	query := map[string]string{
		"markets": markets,
	}

	options := &util.RequestOptions{
		Url:   baseUrl + "/orderbook",
		Query: query,
	}
	err = util.Request(options, &orderbooks)
	return
}

func (client *Client) Accounts() (balances []*types.Balance, err error) {
	token, err := client.Token(nil)
	if err != nil {
		return
	}

	options := &util.RequestOptions{
		Url: baseUrl + "/accounts",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	}
	err = util.Request(options, &balances)
	return
}

func (client *Client) OrderChance(
	marketId string,
) (orderChance types.OrderChance, err error) {
	query := map[string]string{
		"market": marketId,
	}

	token, err := client.Token(query)
	if err != nil {
		return
	}

	options := &util.RequestOptions{
		Url: baseUrl + "/orders/chance",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		Query: query,
	}
	err = util.Request(options, &orderChance)
	return
}

/*
 * get Order list
 * if Market is empty, get all orders
 */
func (client *Client) Orders(
	market string,
	state string,
	page int,
	orderBy string,
) (orders []*types.Order, err error) {

	query := make(map[string]string)

	if len(market) > 0 {
		query["market"] = market
	}

	if len(state) > 0 {
		query["state"] = state
	}

	query["page"] = strconv.Itoa(page)
	query["order_by"] = orderBy

	token, err := client.Token(query)
	if err != nil {
		return
	}

	options := &util.RequestOptions{
		Url: baseUrl + "/orders",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		Query: query,
	}
	err = util.Request(options, &orders)
	return
}

/*
 * get OrderMap for each state list
 */
func (client *Client) OrdersMap(
	market string,
	state string,
	page int,
	orderBy string) (orderMap map[string][]*types.Order, err error) {

	orders, err := client.Orders(market, state, page, orderBy)

	orderMap = make(map[string][]*types.Order)

	if err != nil {
		return
	}

	bidOrders := make([]*types.Order,0)
	askOrders := make([]*types.Order,0)

	if len(orders) > 0 {
		for _, value := range orders {
			if strings.Compare(value.Side, types.ORDERSIDE_BID) == 0 {	//  매수 주문
				bidOrders = append(bidOrders, value)
			} else if strings.Compare(value.Side, types.ORDERSIDE_ASK) == 0 {	// 매도 주문
				askOrders = append(askOrders, value)
			} else {
				// Ignore
			}
		}
		orderMap[types.ORDERSIDE_BID] = bidOrders
		orderMap[types.ORDERSIDE_ASK] = askOrders
	}

	return
}

func (client *Client) OrderByInfo (orderInfo types.OrderInfo) (order *types.Order, err error) {
	return client.Order(orderInfo.Identifier, orderInfo.Side, orderInfo.Market, orderInfo.Price, orderInfo.Volume, orderInfo.OrdType)
}

func (client *Client) Order(
	identifier string,
	side string,
	market string,
	price string,
	volume string,
	ordType string,
) (order *types.Order, err error) {
	query := map[string]string{
		"market":     market,
		"side":       side,
		"volume":     volume,
		"price":      price,
		"ord_type":   ordType,
		"identifier": identifier,
	}

	token, err := client.Token(query)
	if err != nil {
		return
	}

	options := &util.RequestOptions{
		Url:    baseUrl + "/orders",
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		Query: query,
	}

	err = util.Request(options, &order)
	return
}

func (client *Client) CancelOrder(
	uuid string,
) (order *types.Order, err error) {
	query := map[string]string{
		"uuid": uuid,
	}

	token, err := client.Token(query)
	if err != nil {
		return
	}

	options := &util.RequestOptions{
		Url:    baseUrl + "/order",
		Method: "DELETE",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json; charset=utf-8",
		},
		Query: query,
	}

	err = util.Request(options, &order)
	return
}

func NewClient(accessKey, secretKey string) *Client {
	return &Client{accessKey, secretKey}
}

func isValidMinuteCandleUnit(unit int) bool {
	return unit == 1 || unit == 3 || unit == 5 || unit == 10 || unit == 15 ||
		unit == 30 || unit == 60 || unit == 240
}
