package handler

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)

// GetSubscriptionAll ...
// GET /api/v1/subscription/   list subscriptions
func GetSubscriptionAll(w http.ResponseWriter, r *http.Request) {

}

// GetSubscription ...
// GET /api/v1/subscription/:id/   get a subscription
func GetSubscription(w http.ResponseWriter, r *http.Request) {

}

// PostSubscription ...
// POST    /api/v1/subscription/   create a subscription
func PostSubscription(w http.ResponseWriter, r *http.Request) {

}

// UpdateSubscription ...
// PUT /api/v1/subscription/:id/   update a subscription
func PutSubscription(w http.ResponseWriter, r *http.Request) {

}

// DeleteSubscription ...
// DELETE  /api/v1/subscription/:id/   delete a subscription
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {

}

// InitHandlers ...
func InitHandlers(router *mux.Router) {
	apiBase := os.Getenv("API_VERSION_PATH")
	router.HandleFunc(apiBase+"/subscription", GetSubscriptionAll).Methods("GET")
	router.HandleFunc(apiBase+"/subscription/{id}", GetSubscription).Methods("GET")
	router.HandleFunc(apiBase+"/subscription", PostSubscription).Methods("POST")
	router.HandleFunc(apiBase+"/subscription/{id}", PutSubscription).Methods("PUT")
	router.HandleFunc(apiBase+"/subscription/{id}", DeleteSubscription).Methods("DELETE")
}
