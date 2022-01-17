package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jchenriquez/laundromat/store/queries"
	"github.com/jchenriquez/laundromat/types"
)

func getBusinesses(c *gin.Context) {
	query := make(map[string]interface{})
	err := c.BindQuery(query)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	businesses, err := queries.GetBusinesses(db, query)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(businesses)

	if err != nil {
		err = c.Error(err)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func getBusiness(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	business, err := queries.GetBusiness(db, intId)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)

	} else {
		c.Writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(c.Writer)
		err = encoder.Encode(business)
		_ = c.Error(err)
	}
}

func createBusiness(c *gin.Context) {
	var business types.Business
	err := c.Bind(business)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	business, err = queries.CreateBusiness(db, business)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(c.Writer)

	err = encoder.Encode(business)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func businessRouter(r *gin.RouterGroup) {
	businessR := r.Group("/business")
	businessR.GET("/", getBusinesses)
	businessR.GET("/:id", getBusiness)
	businessR.POST("/", createBusiness)
}
