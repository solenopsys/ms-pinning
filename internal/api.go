package internal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
	"ms-pinning/pkg"
	"net/http"
	"strconv"
)

type Api struct {
	Addr        string
	IpfsCluster *pkg.IpfsCluster
	Ipfs        *pkg.IpfsNode
	Data        *Data
}

type BoolResponse struct {
	Result bool `json:"result"`
}

type IdResponse struct {
	Id string `json:"id"`
}

func (api *Api) Start() {
	// http server open port
	r := mux.NewRouter()

	// Define an endpoint to create a new user
	// todo make restfull
	r.HandleFunc("/pin", api.pigGroup).Methods("POST")
	r.HandleFunc("/stat", api.stat).Methods("GET")
	r.HandleFunc("/labels", api.updateLabels).Methods("PUT")
	r.HandleFunc("/stat/pins", api.statPins).Methods("GET")
	r.HandleFunc("/select/pins", api.selectPins).Methods("GET")
	r.HandleFunc("/select/names", api.selectName).Methods("GET")
	r.HandleFunc("/check/pin", api.checkPin).Methods("GET")
	r.HandleFunc("/check/name", api.checkName).Methods("GET")
	r.HandleFunc("/name/create", api.ipnsCreate).Methods("GET")
	r.HandleFunc("/name/update", api.ipnsUpdate).Methods("GET")
	r.HandleFunc("/set/nickname", api.setNickname).Methods("GET")
	r.HandleFunc("/get/pin", api.getPin).Methods("GET")
	r.HandleFunc("/get/name", api.getIpns).Methods("GET")

	// Start the server
	klog.Fatal(http.ListenAndServe(api.Addr, r))
}

type Statistic struct {
	UsersCount int `json:"users_count"`
	PinsCount  int `json:"pins_count"`
	IpnsCount  int `json:"ipns_count"`
}

func (api *Api) getPin(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("cid")
	pin, err := api.Data.GetPin(id)
	if err != nil { // todo change it
		http.Error(w, "Pin not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, pin)
}

func (api *Api) getIpns(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("Authorization")
	userId, err := auth(userKey, api.Data)
	if checkError(err, w) {
		return
	}
	id := r.URL.Query().Get("cid")
	name := r.URL.Query().Get("name")
	var resp *RecordInfo
	if name != "" {
		resp, err = api.Data.GetIpnsByName(name, userId)
	} else {
		resp, err = api.Data.GetIpnsById(id)
	}

	if err != nil { // todo change it
		http.Error(w, "Pin not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, resp)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if checkError(err, w) {
		return
	}
}

func (api *Api) setNickname(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("Authorization")
	nickname := r.URL.Query().Get("value")
	userId, err := auth(userKey, api.Data)
	exists, err := api.Data.CheckNicNameExists(nickname)
	if !exists {
		api.Data.SetNicName(userId, nickname)
	} else {
		http.Error(w, "Nickname already exists", http.StatusBadRequest)
	}

	if checkError(err, w) {
		return
	}
	jsonResponse(w, BoolResponse{Result: true})
}

func (api *Api) checkPin(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("cid")
	result, err := api.Data.PinExists(cid)
	if checkError(err, w) {
		return
	}
	jsonResponse(w, BoolResponse{Result: result})
}

func (api *Api) checkName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("value")
	userKey := r.Header.Get("Authorization")
	userId, err := auth(userKey, api.Data)
	result, err := api.Data.NameExists(name, userId)
	if checkError(err, w) {
		return
	}
	jsonResponse(w, BoolResponse{Result: result})
}

func (api *Api) selectPins(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	value := r.URL.Query().Get("value")

	stat, err := api.Data.SelectPins(name, value)
	if checkError(err, w) {
		return
	}

	jsonResponse(w, stat)
}

func (api *Api) selectName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	value := r.URL.Query().Get("value")

	stat, err := api.Data.SelectIpns(name, value)
	if checkError(err, w) {
		return
	}

	jsonResponse(w, stat)
}

func (api *Api) statPins(w http.ResponseWriter, r *http.Request) {
	stat, err := api.Data.StatByTypes()
	if checkError(err, w) {
		return
	}

	jsonResponse(w, stat)
}

func (api *Api) stat(w http.ResponseWriter, r *http.Request) {
	stat := Statistic{}
	stat.UsersCount, _ = api.Data.GetUsersCount()
	stat.PinsCount, _ = api.Data.GetPinsCount()
	stat.IpnsCount, _ = api.Data.GetIpnsCount()

	jsonResponse(w, stat)
}

func createFullName(userId uint64, name string) string {
	return strconv.FormatUint(userId, 10) + "@" + name
}

func (api *Api) ipnsCreate(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("Authorization")
	userId, err := auth(userKey, api.Data)
	if checkError(err, w) {
		return
	}

	cid := r.URL.Query().Get("cid")
	name := r.URL.Query().Get("name")

	fullName := createFullName(userId, name)

	id, err := api.Ipfs.CreateKey(fullName)
	if checkError(err, w) {
		return
	}

	err = api.Ipfs.Publish(cid, fullName)
	if checkError(err, w) {
		return
	}

	err = api.Data.CreateIpnsRecord(id, userId, cid, name)
	if checkError(err, w) {
		return
	}

	jsonResponse(w, IdResponse{Id: id})
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

	fullName := createFullName(userId, name)
	err = api.Ipfs.Publish(cid, fullName)
	if checkError(err, w) {
		return
	}

	jsonResponse(w, IdResponse{Id: id})
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

func (api *Api) updateLabels(w http.ResponseWriter, r *http.Request) {
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

	for _, pin := range pins.Pins {
		err := api.Data.DeleteLabels(pin.Cid, userId)
		if checkError(err, w) {
			return
		}
		err = api.Data.AddLabels(pin.Labels, pin.Cid, userId, true)
		if checkError(err, w) {
			return
		}
	}

}
