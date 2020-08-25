package apis

import (
	"github.com/go-chi/chi"
	"github.com/jchenriquez/laundromat/store"
	"net/http"
	"strconv"
)

func businessRouter(r chi.Router) {
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		urlQuery := request.URL.Query()
		businesses := make(store.Businesses, 0)

		getCollection(&businesses, urlQuery, writer, request)
	})
	r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		intId, err := strconv.Atoi(id)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		business := store.Business{ID: intId}

		performGet(&business, writer)
	})
	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		performAdd(&store.Business{}, writer, request)
	})
}
