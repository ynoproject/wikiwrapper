package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ynoproject/wikiwrapper/common"
)

func Init() {
	http.HandleFunc("/locations", handleLocations)
	http.HandleFunc("/connections", handleConnections)
	http.HandleFunc("/authors", handleAuthors)
	http.HandleFunc("/maps", handleMaps)
	http.HandleFunc("/vms", handleVendingMachines)

	http.Serve(getListener(), nil)
}

func getListener() net.Listener {
	os.Remove("sockets/wikiwrapper.sock")

	listener, err := net.Listen("unix", "sockets/wikiwrapper.sock")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if err := os.Chmod("sockets/wikiwrapper.sock", 0666); err != nil {
		log.Fatal(err)
		return nil
	}

	return listener
}

func handleLocations(w http.ResponseWriter, r *http.Request) {
	gameParam, ok := r.URL.Query()["game"]
	if !ok || len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	locations, err := common.GetLocations(gameParam[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locationsJson, err := json.Marshal(locations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(locationsJson)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	gameParam, ok := r.URL.Query()["game"]
	if !ok || len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	connections, err := common.GetConnections(gameParam[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	connectionsJson, err := json.Marshal(connections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(connectionsJson)
}

func handleAuthors(w http.ResponseWriter, r *http.Request) {
	gameParam, ok := r.URL.Query()["game"]
	if !ok || len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	if (gameParam[0] != "2kki") && (gameParam[0] != "unevendream") {
		http.Error(w, "game not supported", http.StatusBadRequest)
		return
	}

	authors, err := common.GetAuthors(gameParam[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authorsJson, err := json.Marshal(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(authorsJson)
}

func handleMaps(w http.ResponseWriter, r *http.Request) {
	gameParam, ok := r.URL.Query()["game"]
	if !ok || len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	locationParam, ok := r.URL.Query()["location"]
	if !ok || len(locationParam) == 0 {
		http.Error(w, "location not specified", http.StatusBadRequest)
		return
	}

	maps, err := common.GetMaps(gameParam[0], locationParam[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mapsJson, err := json.Marshal(maps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(mapsJson)
}

func handleVendingMachines(w http.ResponseWriter, r *http.Request) {
	gameParam, ok := r.URL.Query()["game"]
	if !ok || len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	vms, err := common.GetVendingMachines(gameParam[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vmsJson, err := json.Marshal(vms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(vmsJson)
}
