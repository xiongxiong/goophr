package main

import (
	"fmt"
	"goophr/librarian/api"
	"goophr/librarian/common"
	"net/http"
	"os"
)

func main() {
	common.Log("Adding API handlers...")
	http.HandleFunc("/api/index", api.IndexHandler)
	http.HandleFunc("/api/query", api.QueryHandler)

	common.Log("Starting index...")
	api.StartIndexSystem()

	port := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	common.Log("Starting Goophr Librarian server on port :" + port + "...")
	http.ListenAndServe(port, nil)
}
