package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudrimmers/imt2681-assignment3/lib/dialogFlow"
	"github.com/gorilla/mux"
)

const root = "/rimbot/"

var err error

//Rimbot - TODO
func Rimbot(w http.ResponseWriter, r *http.Request) {
	//obtain slacks' webhook
	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := r.Form

	//convert webhook values into new outgoing message
	base, target, amount, code := dialogFlow.Query(form.Get("text"))
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}
	test := fmt.Sprintf("Response got:\n%v\t%v\t%v", base, target, amount)

	fmt.Println("test: ", test)
	w.Write([]byte(test))

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rimbot/", Rimbot).Methods(http.MethodPost)
	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
