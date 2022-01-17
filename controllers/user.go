package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jchenriquez/laundromat/store/queries"
	"github.com/jchenriquez/laundromat/types"
)

func getUserProfile(c *gin.Context) {
	id := c.Param("id")

	user, err := queries.GetUserById(db, id)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = ""
	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(user)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()
}

func updateUserProfile(c *gin.Context) {

	var user types.User
	email := c.Param("email")
	err := c.BindJSON(user)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	user.Email = email
	user, err = queries.UpdateUser(db, user)
	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = ""
	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(user)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()
}

func userRouter(r *gin.RouterGroup) {
	userR := r.Group("/user")
	userR.GET("/:id", getUserProfile)
	userR.PUT("/:id", updateUserProfile)
}
