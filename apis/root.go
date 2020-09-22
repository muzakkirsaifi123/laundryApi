package apis

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jchenriquez/laundromat/store"
	config "github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
)

var clientDB *store.Client

const (
	SECOND = 1
	MINUTE = 60 * SECOND
)

func fillModelWithData(model store.Model) {
	err := model.Get(clientDB)
	if err != nil {
		panic(err)
	}
}

func fillCollectionWithData(collection store.Collection, query store.Query) {
	err := collection.Get(clientDB, query)
	if err != nil {
		panic(err)
	}
}

func writeModelData(model store.Model, writer http.ResponseWriter) {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(model)

	if err != nil {
		panic(err)
	}
}

func writeModelToStorage(model store.Model) {
	err := model.Create(clientDB)

	if err != nil {
		panic(err)
	}
}

func writeCollectionData(collection store.Collection, writer http.ResponseWriter) {
	bufioWriter := bufio.NewWriter(writer)
	encoder := json.NewEncoder(bufioWriter)
	err := encoder.Encode(collection)

	if err != nil {
		panic(err)
	}

	err = bufioWriter.Flush()

	if err != nil {
		panic(err)
	}
}

func setClientDB() {
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

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(writer, err.(error).Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func jsonDecodeRequestBody(model interface{}, request *http.Request) {
	jsonDecoder := json.NewDecoder(request.Body)
	var err error

	err = jsonDecoder.Decode(model)

	if err != nil {
		panic(err)
	}
}

func AddApis(router chi.Router) {
	setClientDB()

	router.Use(errorHandler)

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
		r.Route("/business", businessRouter)
		r.Route("/user", userRouter)
		r.Route("/order", orderRouter)
		r.Route("/session", sessionRouter)
	})

}
