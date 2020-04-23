package controller

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type Order struct {
	Number          int64   `json:"number"`
	CreateTimestamp int64   `json:"create_time"`
	Symbol          string  `json:"symbol"`
	Type            int     `json:"type"`
	Lot             float64 `json:"lot"`
	Price           float64 `json:"price"`
	StopLoss        float64 `json:"stop_loss"`
	TakeProfit      float64 `json:"take_profit"`
	Open            bool    `json:"open"`
}

func (o *Order) Equal(in *Order) bool {
	return o.Number == in.Number && o.CreateTimestamp == in.CreateTimestamp && o.Symbol == in.Symbol &&
		o.Type == in.Type && o.Lot == in.Lot && o.Price == in.Price && o.StopLoss == in.StopLoss && o.TakeProfit == in.TakeProfit && o.Open == in.Open
}

type ByNumber []*Order

func (a ByNumber) Len() int {
	return len(a)
}

func (a ByNumber) Less(i, j int) bool {
	return a[i].Number < a[j].Number
}

func (a ByNumber) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type PostOrdersBody struct {
	Orders []*Order `json:"orders"`
}

type OrdersController struct {
	Orders []*Order
}

func (ctr *OrdersController) PostController(c *gin.Context) {
	var req PostOrdersBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sort.Sort(ByNumber(req.Orders))
	var newOrder, closedOrder, modifiedOrder []*Order

	for _, order := range ctr.Orders {
		found := false
		for _, orderNew := range req.Orders {
			if order.Number == orderNew.Number {
				found = true
				if !order.Equal(orderNew) {
					modifiedOrder = append(modifiedOrder, orderNew)
				}
				break
			} else if order.Number < orderNew.Number {
				break
			}
		}
		if found == false {
			closedOrder = append(closedOrder, order)
		}
	}

	for _, orderNew := range req.Orders {
		found := false
		for _, order := range ctr.Orders {
			if orderNew.Number == order.Number {
				found = true
				break
			} else if orderNew.Number < order.Number {
				break
			}
		}
		if found == false {
			newOrder = append(newOrder, orderNew)
		}
	}

	fmt.Println("new")
	for _, order := range newOrder {
		fmt.Println(order.Number)
	}
	fmt.Println("closed")
	for _, order := range closedOrder {
		fmt.Println(order.Number)
	}
	fmt.Println("modified")
	for _, order := range modifiedOrder {
		fmt.Println(order.Number)
	}

	ctr.Orders = req.Orders

	c.Status(http.StatusOK)
}
