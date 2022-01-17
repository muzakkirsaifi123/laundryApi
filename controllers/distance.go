package controllers

//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/go-chi/chi"
// 	config "github.com/spf13/viper"
// 	"googlemaps.github.io/maps"
// 	"net/http"
// 	"strconv"
//
// 	"github.com/jchenriquez/laundromat/store"
// )
//
// func transportRouter(r chi.Router) {
// 	r.Group(func(r chi.Router) {
// 		r.Use(Auth(clientDB))
// 		r.Get("/distance/{uid}/{businessid}", func(writer http.ResponseWriter, request *http.Request) {
// 			var client *maps.Client
// 			fromUserId := chi.URLParam(request, "uid")
// 			toBusinessId := chi.URLParam(request, "businessid")
// 			distanceApiKey := config.GetString("distance_api_key")
//
// 			client, err := maps.NewClient(maps.WithAPIKey(distanceApiKey))
//
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			uid, err := strconv.Atoi(fromUserId)
//
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			businessId, err := strconv.Atoi(toBusinessId)
//
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			user := store.User{UID: uid}
// 			business := store.Business{ID: businessId}
//
// 			fillModelWithData(&user)
// 			fillModelWithData(&business)
//
// 			r := &maps.DistanceMatrixRequest{
// 				Mode:         "Driving",
// 				Origins:      []string{fmt.Sprintf("%s, %s, %s, %s, %s", user.AddressLine1, user.AddressLine2, user.City, user.State, user.ZipCode)},
// 				Destinations: []string{business.Address},
// 			}
//
// 			distanceResponse, err := client.DistanceMatrix(request.Context(), r)
//
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			distanceInfo := distanceResponse.Rows[0].Elements[0]
//
// 			encoder := json.NewEncoder(writer)
// 			err = encoder.Encode(distanceInfo)
//
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			writer.WriteHeader(http.StatusOK)
// 		})
// 	})
// }
