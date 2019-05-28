package util

import (
	"fmt"
	"github.com/jekeun/upbit-go/types"
)

func PrintBalance(balances []*types.Balance) {
	for _, value := range balances {
		fmt.Printf("%+v\n", value)
	}
}

func PrintOrdersMap(ordersMap map[string][]*types.Order) {
	if ordersMap == nil || len(ordersMap) <= 0 {
		return
	}

	bidOrders := ordersMap[types.ORDERSIDE_BID]

	fmt.Println("Bid Orders")
	for _, value := range bidOrders {
		fmt.Printf("%+v\n", value)
	}

	askOrders := ordersMap[types.ORDERSIDE_ASK]

	fmt.Println("Ask Orders")
	for _, value := range askOrders {
		fmt.Printf("%+v\n", value)
	}
}