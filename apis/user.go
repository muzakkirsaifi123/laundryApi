package apis

import (
	"github.com/go-chi/chi"
	"github.com/jchenriquez/laundromat/store"
	"net/http"
	"strconv"
)

func userRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(Auth(clientDB))
		r.Get("/{id:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
			id := chi.URLParam(request, "id")
			intId, err := strconv.Atoi(id)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			writeModelData(&store.User{UID: intId}, writer)
		})
	})
	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		user := store.User{}
		jsonDecodeRequestBody(&user, request)
		writeModelToStorage(&user)
	})
}
