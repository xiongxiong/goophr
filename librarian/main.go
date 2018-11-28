package main

import (
	"goophr/librarian/api"
	"goophr/librarian/common"
	"net/http"
)

func main() {
	common.Log("Adding API handlers...")
	http.HandleFunc("api/index", api.IndexHandler)
	http.HandleFunc("api/query", api.QueryHandler)

	common.Log("Starting index...")
	api.StartIndexSystem()

	common.Log("Starting Goophr Librarian server on port :9090...")
	http.ListenAndServe(":9090", nil)
}
