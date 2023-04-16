package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
	"net/http"
)

type Api struct {
	addr string
	ipfs *IpfsCluster
	data *Data
}

func (api *Api) Start() {
	// http server open port
	r := mux.NewRouter()

	// Define an endpoint to create a new user
	r.HandleFunc("/pin", api.pigGroup).Methods("POST")

	// Start the server
	klog.Fatal(http.ListenAndServe(api.addr, r))
}

func (api *Api) pigGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userKey := r.Header.Get("Authorization")
	var pins IpfsPingGroup
	err := json.NewDecoder(r.Body).Decode(&pins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId int64
	user := api.data.GetUserById(userKey)
	if user == nil {
		userId, err = api.data.AddUser(userKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		userId = user.id
	}

	group := api.ipfs.PinGroup(pins)

	for _, pin := range group {
		pinId := pin.Cid.String()
		err := api.data.AddPin(pinId, userId, 0)

		if err != nil {
			klog.Error("Pin save in db error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			klog.Info("Pin save in db: ", pinId)
		}
	}

	// Return the new user as JSON
	json.NewEncoder(w).Encode(group)
}
