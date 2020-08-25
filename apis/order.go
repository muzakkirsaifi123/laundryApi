package apis

import (
	"github.com/go-chi/chi"
	"github.com/jchenriquez/laundromat/store"
	"net/http"
	"strconv"
)

func orderRouter(r chi.Router) {
	r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		intId, err := strconv.Atoi(id)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		performGet(&store.Order{ID: intId}, writer)
	})
	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		performAdd(&store.Order{}, writer, request)
	})

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		orders := make(store.Orders, 0)

		getCollection(&orders, urlQuery, writer, request)
	})
}
