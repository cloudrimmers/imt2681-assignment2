package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Arxcis/imt2681-assignment2/lib/mytypes"
)

// Example: router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", gitRepositoryHandler)
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

// PostWebhook ...
// POST    /api/v1/subscription/   create a subscription
func PostWebhook(w http.ResponseWriter, r *http.Request) {

	webhook := &mytypes.WebhookIn{}
	json.NewDecoder(r.Body).Decode(webhook)
	fmt.Println(webhook)
	strrandom := strconv.Itoa(rand.Int())
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strrandom))
}
