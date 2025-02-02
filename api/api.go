package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ynoproject/wikiwrapper/common"
	"github.com/ynoproject/wikiwrapper/setup"
)

func Init() {
	corsConfig, err := setup.LoadCorsConfig("cors_config.yml")
	if err != nil {
		log.Fatalf("Error loading CORS config: %v", err)
	}

	wikiConfig, err := setup.LoadWikiConfig("wiki_config.yml")
	if err != nil {
		log.Println(wikiConfig.Games)
		log.Fatalf("Error loading wiki config: %v", err)
	}

	http.HandleFunc("/locations", handleLocations)
	http.HandleFunc("/connections", handleConnections)
	http.HandleFunc("/authors", handleAuthors)
	http.HandleFunc("/maps", handleMaps)
	http.HandleFunc("/vms", handleVendingMachines)
	http.HandleFunc("/images", handleImages)

	configMiddleware := setup.WikiConfigHandlerMiddleware(wikiConfig)
	corsHandler := setup.CorsHandlerMiddleware(corsConfig)
	handler := configMiddleware(corsHandler)

	http.Serve(getListener(), handler)
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
	config := r.Context().Value(setup.ConfigKey).(setup.WikiConfig)
	gameParam := r.URL.Query().Get("game")
	if gameParam == "" {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	protagParam := r.URL.Query().Get("protag")
	gameParams := common.GameParams{GameCode: gameParam}
	if protagParam != "" {
		gameParams.Protag = protagParam
	}

	continueKeyParam := r.URL.Query().Get("continueKey")
	if continueKeyParam != "" {
		gameParams.ContinueKey = continueKeyParam
	}

	locations, err := common.GetLocations(gameParams, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locationsJson, err := json.Marshal(locations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(locationsJson)
}

func handleImages(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(setup.ConfigKey).(setup.WikiConfig)
	gameParam := r.URL.Query().Get("game")
	if gameParam == "" {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	gameParams := common.GameParams{GameCode: gameParam}
	continueKeyParam := r.URL.Query().Get("continueKey")
	if continueKeyParam != "" {
		gameParams.ContinueKey = continueKeyParam
	}

	images, err := common.GetImages(gameParams, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imagesJson, err := json.Marshal(images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(imagesJson)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(setup.ConfigKey).(setup.WikiConfig)
	gameParam := r.URL.Query().Get("game")
	if gameParam == "" {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	protagParam := r.URL.Query().Get("protag")
	gameParams := common.GameParams{GameCode: gameParam}
	if protagParam != "" {
		gameParams.Protag = protagParam
	}

	continueKeyParam := r.URL.Query().Get("continueKey")
	if continueKeyParam != "" {
		gameParams.ContinueKey = continueKeyParam
	}

	connections, err := common.GetConnections(gameParams, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	connectionsJson, err := json.Marshal(connections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(connectionsJson)
}

func handleAuthors(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(setup.ConfigKey).(setup.WikiConfig)
	gameParam := r.URL.Query().Get("game")
	if len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	if (gameParam != "2kki") && (gameParam != "unevendream") && (gameParam != "unconscious") {
		http.Error(w, "game not supported", http.StatusBadRequest)
		return
	}

	authors, err := common.GetAuthors(gameParam, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authorsJson, err := json.Marshal(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(authorsJson)
}

func handleMaps(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(setup.ConfigKey).(setup.WikiConfig)
	gameParam := r.URL.Query().Get("game")
	if len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	locationParam := r.URL.Query().Get("location")
	if len(locationParam) == 0 {
		http.Error(w, "location not specified", http.StatusBadRequest)
		return
	}

	maps, err := common.GetMaps(gameParam, locationParam, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mapsJson, err := json.Marshal(maps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(mapsJson)
}

func handleVendingMachines(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(setup.ConfigKey).(setup.WikiConfig)
	gameParam := r.URL.Query().Get("game")
	if len(gameParam) == 0 {
		http.Error(w, "game not specified", http.StatusBadRequest)
		return
	}

	vms, err := common.GetVendingMachines(gameParam, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vmsJson, err := json.Marshal(vms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(vmsJson)
}
