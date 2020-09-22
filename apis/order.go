package apis

import (
	"github.com/go-chi/chi"
	"github.com/jchenriquez/laundromat/store"
	"net/http"
	"strconv"
)

func orderRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(Auth(clientDB))
		r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
			id := chi.URLParam(request, "id")
			intId, err := strconv.Atoi(id)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			order := &store.Order{ID: intId}
			fillModelWithData(order)
			writeModelData(order, writer)
		})
		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			order := store.Order{}
			jsonDecodeRequestBody(&order, request)
			writeModelToStorage(&order)
		})

		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			urlQuery := request.URL.Query()
			orders := make(store.Orders, 0)

			fillCollectionWithData(&orders, urlQuery)
			writeCollectionData(&orders, writer)
		})
	})
}
