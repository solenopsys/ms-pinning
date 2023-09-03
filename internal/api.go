package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
	"ms-pinning/pkg"
	"net/http"
)

type Api struct {
	Addr        string
	IpfsCluster *pkg.IpfsCluster
	Ipfs        *pkg.IpfsNode
	Data        *Data
}

func (api *Api) Start() {
	// http server open port
	r := mux.NewRouter()

	// Define an endpoint to create a new user
	r.HandleFunc("/pin", api.pigGroup).Methods("POST")
	r.HandleFunc("/stat", api.stat).Methods("GET")
	r.HandleFunc("/ipns/create", api.ipnsCreate).Methods("GET")
	r.HandleFunc("/ipns/update", api.ipnsUpdate).Methods("GET")

	// Start the server
	klog.Fatal(http.ListenAndServe(api.Addr, r))
}

type Statistic struct {
	UsersCount int `json:"users_count"`
	PinsCount  int `json:"pins_count"`
	IpnsCount  int `json:"ipns_count"`
}

func (api *Api) stat(w http.ResponseWriter, r *http.Request) {
	stat := Statistic{}
	stat.UsersCount, _ = api.Data.GetUsersCount()
	stat.PinsCount, _ = api.Data.GetPinsCount()
	stat.IpnsCount, _ = api.Data.GetIpnsCount()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(stat)
	if checkError(err, w) {
		return
	}
}

func (api *Api) ipnsCreate(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("Authorization")
	userId, err := auth(userKey, api.Data)
	if checkError(err, w) {
		return
	}

	cid := r.URL.Query().Get("cid")
	name := r.URL.Query().Get("name")

	id, err := api.Ipfs.CreateKey(name)
	if checkError(err, w) {
		return
	}

	err = api.Ipfs.Publish(cid, name)
	if checkError(err, w) {
		return
	}
	err = api.Data.CreateIpnsRecord(id, userId, cid, name)
	if checkError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/text")
	_, err = w.Write([]byte(id))
	if checkError(err, w) {
		return
	}

}

func (api *Api) ipnsUpdate(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("Authorization")
	userId, err := auth(userKey, api.Data)
	if checkError(err, w) {
		return
	}

	cid := r.URL.Query().Get("cid")
	name := r.URL.Query().Get("name")

	id, err := api.Data.ChangeIpnsRecord(name, cid, userId)
	if checkError(err, w) {
		return
	}
	err = api.Ipfs.Publish(cid, name)
	if checkError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/text")
	_, err = w.Write([]byte(id))
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

	group, err := api.IpfsCluster.PinGroup(pins)
	if checkError(err, w) {
		return
	}

	err = api.Data.SavePins(pins.Pins, group, userId)
	if checkError(err, w) {
		return
	}
}
