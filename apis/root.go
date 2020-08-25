package apis

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jchenriquez/laundromat/apis/lib"
	"github.com/jchenriquez/laundromat/store"
	config "github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var clientDB *store.Client

const (
	SECOND = 1
	MINUTE = 60 * SECOND
)

func performModelGet(writer io.Writer, model store.Model) error {
	err := model.Get(clientDB)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(writer)
	return encoder.Encode(model)
}

func performGet(model store.Model, writer http.ResponseWriter) {
	err := performModelGet(writer, model)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func performAdd(model store.Model, writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(model)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.Create(clientDB)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getCollection(collection store.Collection, query store.Query, writer http.ResponseWriter, request *http.Request) {
	err := collection.Get(clientDB, query)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	bufioWriter := bufio.NewWriter(writer)
	encoder := json.NewEncoder(bufioWriter)
	err = encoder.Encode(collection)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	err = bufioWriter.Flush()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func start() {
	dbName := config.Get("database_name").(string)
	dbHostname := config.GetString("database_hostname")
	dbPort := config.GetString("database_port")
	dbUsername := config.GetString("database_username")
	dbPassword := config.GetString("database_password")

	port, err := strconv.Atoi(dbPort)

	if err != nil {
		log.Fatal(err)
	}

	clientDB = store.New(context.Background(), dbUsername, dbHostname, dbPassword, dbName, port, false)
}

func reAuthorize(usr store.User) (store.Session, error) {
	accessToken, expirationStr, err := lib.GenerateAccessToken(usr)

	if err != nil {
		return store.Session{}, err
	}

	refreshToken, _, err := lib.GenerateRefreshToken(usr)

	if err != nil {
		return store.Session{}, err
	}

	return store.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiration:   expirationStr,
	}, nil
}

func AddApis(router chi.Router) {
	start()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			err := clientDB.Open()

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			defer func() {
				err = clientDB.Close()

				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(writer, request)
		})
	})
	router.Group(func(r chi.Router) {
		r.Use(lib.Auth(clientDB))
		r.Route("/business", businessRouter)
		r.Route("/user", userRouter)
		r.Route("/order", orderRouter)
		r.Post("/user/logOff", func(writer http.ResponseWriter, request *http.Request) {
			session := store.Session{}
			decoder := json.NewDecoder(request.Body)
			err := decoder.Decode(&session)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			err = session.Delete(clientDB)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(lib.RefreshAuth(clientDB))
		r.Post("/user/refresh", func(writer http.ResponseWriter, request *http.Request) {
			_, claims, err := jwtauth.FromContext(request.Context())

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			usr := store.User{UID: int(claims["UID"].(float64))}
			tokenPayload, err := reAuthorize(usr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			session := store.Session{UID: usr.UID, AccessToken: tokenPayload.AccessToken, RefreshToken: tokenPayload.RefreshToken}
			err = session.Create(clientDB)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			encoder := json.NewEncoder(writer)
			err = encoder.Encode(tokenPayload)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			writer.WriteHeader(http.StatusOK)
		})
	})

	router.Post("/user/login", func(writer http.ResponseWriter, request *http.Request) {
		jsonDecoder := json.NewDecoder(request.Body)
		user := store.User{}
		err := jsonDecoder.Decode(&user)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		user.UID = -1
		// Plain text password from payload
		password := user.Password

		err = user.Get(clientDB)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		if user.UID == -1 {
			http.Error(writer, "user does not exist", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		user.Password = ""
		tokenPayload, err := reAuthorize(user)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		session := store.Session{UID: user.UID, AccessToken: tokenPayload.AccessToken, RefreshToken: tokenPayload.RefreshToken}
		err = session.Create(clientDB)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonEncoder := json.NewEncoder(writer)
		err = jsonEncoder.Encode(tokenPayload)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

}
