package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
	"net/http"
)

type Api struct {
	Addr string
	Ipfs *IpfsCluster
	Data *Data
}

func (api *Api) Start() {
	// http server open port
	r := mux.NewRouter()

	// Define an endpoint to create a new user
	r.HandleFunc("/pin", api.pigGroup).Methods("POST")
	r.HandleFunc("/stat", api.stat).Methods("GET")

	// Start the server
	klog.Fatal(http.ListenAndServe(api.Addr, r))
}

func (api *Api) stat(w http.ResponseWriter, r *http.Request) {

}

func checkError(err error, w http.ResponseWriter) bool {
	isErr := err != nil
	if isErr {
		klog.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return isErr
}

func (api *Api) pigGroup(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("Authorization")
	userId, err := auth(userKey, api.Data)
	if checkError(err, w) {
		return
	}

	var pins IpfsPinsGroup
	err = json.NewDecoder(r.Body).Decode(&pins)
	if checkError(err, w) {
		return
	}

	group, err := api.Ipfs.PinGroup(pins)
	if checkError(err, w) {
		return
	}

	err = api.Data.SavePins(pins.Pins, group, userId)
	if checkError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(group)
	if checkError(err, w) {
		return
	}
}
