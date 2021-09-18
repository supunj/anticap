// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package api

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterRoutes - Register routes
func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Swagger documentation
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Routes consist of a path and a handler function.
	router.HandleFunc("/api/v1/register/{node}", Register).Methods("POST")
	router.HandleFunc("/api/v1/verify/{nodeid}/{vcode}", Verify).Methods("PUT")
	router.HandleFunc("/api/v1/subscribe/{nodeid}/{channel}/{as}", Subscribe).Methods("PUT")
	router.HandleFunc("/api/v1/sendrequest/{nodeid}/{channel}/{request}", SendRequest).Methods("POST")
	router.HandleFunc("/api/v1/cancelrequest/{nodeid}/{requestid}/{reason}", CancelRequest).Methods("PUT")

	// WebSocket endpoint to receive requests
	router.HandleFunc("/api/v1/receiverequests", ServeRequests)

	// Provider makes an offer
	router.HandleFunc("/api/v1/makeoffer/{nodeid}/{requestid}/{offer}", MakeOffer).Methods("POST")

	// WebSocket endpoint to listen and negotiate the offer
	router.HandleFunc("/api/v1/negotiateoffer/{offerid}", NegotiateOffer).Methods("POST")

	// WebSocket endpoint to receive responces from providers
	router.HandleFunc("/api/v1/receiveresponces/{requestid}", ReceiveResponse).Methods("GET")

	// Accept or reject offer
	router.HandleFunc("/api/v1/acceptoffer/{nodeid}/{offerid}", AcceptOffer).Methods("PUT")
	router.HandleFunc("/api/v1/rejectoffer/{nodeid}/{offerid}", RejectOffer).Methods("PUT")

	router.HandleFunc("/api/v1/complete", Complete).Methods("PUT")
	router.HandleFunc("/api/v1/feedback", Feedback).Methods("POST")

	// Regular calls
	router.HandleFunc("/api/v1/updatelocation/{nodeid}/{lat}/{lon}", UpdateLocation).Methods("PUT")
	router.HandleFunc("/api/v1/refreshkeys/{nodeid}/{oldpvtkey}", RefreshKeys).Methods("POST")

	return router
}
