package controller

import (
	"net/http"

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

type PostOrdersBody struct {
	Orders []*Order `json:"orders"`
}

type OrdersController struct{}

func (ctr *OrdersController) PostController(c *gin.Context) {
	var req PostOrdersBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
