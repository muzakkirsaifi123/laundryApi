package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jchenriquez/laundromat/store/queries"
	"github.com/jchenriquez/laundromat/types"
)

func getOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := queries.GetOrder(db, id)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(order)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}

	c.Writer.Flush()
}

func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var order types.Order
	err := c.BindJSON(order)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	order.ID = id
	order, err = queries.UpdateOrder(db, order)
	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(order)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}
	c.Writer.Flush()
}

func createOrder(c *gin.Context) {
	var order types.Order
	err := c.BindJSON(order)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	order.Status = "Awaiting Pickup"
	order, err = queries.CreateOrder(db, order)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Writer.Flush()
}

func getOrders(c *gin.Context) {
	query := make(map[string]interface{})
	err := c.BindQuery(query)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	var orders []types.Order
	orders, err = queries.GetOrders(db, query)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(orders)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}
	c.Writer.Flush()
}

func orderRouter(r *gin.RouterGroup) {
	ordersR := r.Group("/order")
	ordersR.GET("/:id", getOrder)
	ordersR.PUT("/:id", updateOrder)
	ordersR.POST("/", createOrder)
	ordersR.GET("/", getOrders)
}
