package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	addr string
	ipfs *IpfsCluster
	db   *sql.DB
	data *Data
}

func (api *Api) Start() {
	// http server open port
	r := mux.NewRouter()

	// Define an endpoint to create a new user
	r.HandleFunc("/pin", api.pigGroup)

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

	for _, cid := range pins.Cids {
		api.data.AddPin(cid, userId, 10)
	}

	group := api.ipfs.PinGroup(pins)

	// Return the new user as JSON
	json.NewEncoder(w).Encode(group)
}
