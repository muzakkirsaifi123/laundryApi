package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jchenriquez/laundromat/auth"
	"github.com/jchenriquez/laundromat/store/queries"
	"github.com/jchenriquez/laundromat/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/jchenriquez/laundromat/store"
)

const (
	SESSION_NAME = "__Secure_ID"
)

var db *store.Client

func encryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	sessionManager, err := auth.NewSessionManager(c, db)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := queries.GetUserByEmail(db, email)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, "user does not exist", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	session := sessionManager.CreateSession(SESSION_NAME, user.ID)
	err = sessionManager.SaveSession(session)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(user)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}
	c.Writer.Flush()
}

func logOff(c *gin.Context) {
	sessionManager, err := auth.NewSessionManager(c, db)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessionManager.DeleteSession(SESSION_NAME)
	if err != nil {
		_ = c.Error(err)
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func signup(c *gin.Context) {
	var user types.User

	err := c.BindJSON(&user)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusPreconditionFailed)
		return
	}

	if len(user.Email) == 0 || len(user.Password) == 0 {
		http.Error(c.Writer, "email and password must be provided", http.StatusPreconditionFailed)
		return
	}

	user.Password, err = encryptPassword(user.Password)

	if err != nil {
		err = c.Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err = queries.CreateUser(db, user)

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

	c.Writer.Flush()
}

func authorizer(c *gin.Context) {
	sessionManager, err := auth.NewSessionManager(c, db)

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session, err := sessionManager.GetSession(SESSION_NAME)

	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if !session.Valid() {
		_ = sessionManager.DeleteSession(SESSION_NAME)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = sessionManager.RefreshSession(session)

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Next()
}

func storeConnector(c *gin.Context) {
	err := db.Open()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Next()
	defer func() {
		err = db.Close()
		if err != nil {
			_ = c.Error(err)
		}
	}()
}

func AddControllers(router gin.IRouter, dbClient *store.Client) {
	db = dbClient

	router.Use(storeConnector)
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/logoff", logOff)
	securedGroup := router.Group("/admin")
	securedGroup.Use(authorizer)
	businessRouter(securedGroup)
	userRouter(securedGroup)
	orderRouter(securedGroup)
}
