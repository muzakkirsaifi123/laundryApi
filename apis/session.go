package apis

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jchenriquez/laundromat/store"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func createSession(usr store.User) (store.Session, error) {
	accessToken, expirationStr, err := GenerateAccessToken(usr)

	if err != nil {
		return store.Session{}, err
	}

	refreshToken, _, err := GenerateRefreshToken(usr)

	if err != nil {
		return store.Session{}, err
	}

	return store.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiration:   expirationStr,
	}, nil
}

func refresh(writer http.ResponseWriter, request *http.Request) {
	_, claims, err := jwtauth.FromContext(request.Context())

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	usr := store.User{UID: int(claims["UID"].(float64))}
	tokenPayload, err := createSession(usr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	session := store.Session{UID: usr.UID, AccessToken: tokenPayload.AccessToken, RefreshToken: tokenPayload.RefreshToken}
	writeModelToStorage(&session)
	writeModelData(&session, writer)
}

func login(writer http.ResponseWriter, request *http.Request) {
	user := store.User{}
	jsonDecodeRequestBody(&user, request)
	user.UID = -1
	// Plain text password from payload
	password := user.Password

	fillModelWithData(&user)

	if user.UID == -1 {
		http.Error(writer, "user does not exist", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	user.Password = ""
	tokenPayload, err := createSession(user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	session := store.Session{UID: user.UID, AccessToken: tokenPayload.AccessToken, RefreshToken: tokenPayload.RefreshToken}
	writeModelToStorage(&session)
	writeModelData(&session, writer)
}

func logOff(writer http.ResponseWriter, request *http.Request) {
	session := store.Session{}
	jsonDecodeRequestBody(&session, request)

	err := session.Delete(clientDB)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sessionRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(RefreshAuth(clientDB))
		r.Post("/refresh", refresh)
	})
	r.Post("/login", login)
	r.Group(func(r chi.Router) {
		r.Use(Auth(clientDB))
		r.Post("/logOff", logOff)
	})
}
