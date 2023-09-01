package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
	"ms-pinning/pkg"
	"net/http"
)

type Api struct {
	Addr string
	Ipfs *pkg.IpfsCluster
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

type Statistic struct {
	UsersCount int `json:"users_count"`
	PinsCount  int `json:"pins_count"`
}

func (api *Api) stat(w http.ResponseWriter, r *http.Request) {
	stat := Statistic{}
	stat.UsersCount, _ = api.Data.GetUsersCount()
	stat.PinsCount, _ = api.Data.GetPinsCount()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(stat)
	if checkError(err, w) {
		return
	}
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

	var pins pkg.IpfsPinsGroup
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
}
