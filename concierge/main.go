package main

import (
	"fmt"
	"net/http"
	"os"

	"goophr/concierge/api"
	"goophr/concierge/common"
)

func main() {
	common.Log("Adding API handlers...")
	http.HandleFunc("/api/feeder", api.FeedHandler)
	http.HandleFunc("/api/query", api.QueryHandler)

	common.Log("Starting feeder...")
	api.StartFeederSystem()

	port := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	common.Log("Starting Goophr Concierge server on port :" + port + "...")
	http.ListenAndServe(port, nil)
}
