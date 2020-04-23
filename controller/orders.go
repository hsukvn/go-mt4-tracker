package controller

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hsukvn/go-mt4-tracker/util"
)

type Order struct {
	Number          int64   `json:"number"`
	CreateTimestamp int64   `json:"create_time"`
	Symbol          string  `json:"symbol"`
	Type            int     `json:"type"`
	Lot             float64 `json:"lot"`
	Digit           int     `json:"digit"`
	Price           float64 `json:"price"`
	StopLoss        float64 `json:"stop_loss"`
	TakeProfit      float64 `json:"take_profit"`
	Open            bool    `json:"open"`
}

func OrderTypeString(t int) string {
	m := map[int]string{
		0: "BUY",
		1: "SELL",
		2: "BUYLIMIT",
		3: "SELLLIMIT",
		4: "BUYSTOP",
		5: "SELLSTOP",
	}
	return m[t]
}

func (o *Order) Equal(in *Order) bool {
	return o.Number == in.Number && o.CreateTimestamp == in.CreateTimestamp && o.Symbol == in.Symbol &&
		o.Type == in.Type && o.Lot == in.Lot && o.Digit == in.Digit && o.Price == in.Price && o.StopLoss == in.StopLoss && o.TakeProfit == in.TakeProfit && o.Open == in.Open
}

func (o *Order) String() string {
	// Number/Type/Symbol/Lot/Price/SL/TP
	return strconv.FormatInt(o.Number, 10) + "/" + OrderTypeString(o.Type) + "/" + o.Symbol + "/" + strconv.FormatFloat(o.Lot, 'f', 2, 64) + "/" + strconv.FormatFloat(o.Price, 'f', o.Digit, 64) +
		"/" + strconv.FormatFloat(o.StopLoss, 'f', o.Digit, 64) + "/" + strconv.FormatFloat(o.TakeProfit, 'f', o.Digit, 64)
}

func getOrdersString(os []*Order) string {
	s := ""
	for _, o := range os {
		if len(s) == 0 {
			s = o.String()
		} else {
			s = s + "\n" + o.String()
		}
	}
	return s
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
	var newOrders, closedOrders, modifiedOrders []*Order

	for _, order := range ctr.Orders {
		found := false
		for _, orderNew := range req.Orders {
			if order.Number == orderNew.Number {
				found = true
				if !order.Equal(orderNew) {
					modifiedOrders = append(modifiedOrders, orderNew)
				}
				break
			} else if order.Number < orderNew.Number {
				break
			}
		}
		if found == false {
			closedOrders = append(closedOrders, order)
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
			newOrders = append(newOrders, orderNew)
		}
	}

	if len(newOrders) > 0 {
		msg := getOrdersString(newOrders)
		msg = "new orders\n" + msg

		util.SendNotify(msg)
	}

	if len(closedOrders) > 0 {
		msg := getOrdersString(closedOrders)
		msg = "closed orders\n" + msg

		util.SendNotify(msg)
	}

	if len(modifiedOrders) > 0 {
		msg := getOrdersString(modifiedOrders)
		msg = "modified orders\n" + msg

		util.SendNotify(msg)
	}

	ctr.Orders = req.Orders

	c.Status(http.StatusOK)
}
