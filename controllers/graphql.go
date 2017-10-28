package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bitbucket.org/mlsdatatools/retsetl/config"

	"bitbucket.org/mlsdatatools/retsetl/graph"
	"bitbucket.org/mlsdatatools/retsetl/helper"
)

//GetGraphiQL Query interface that can be accessed from http://server:port/grapihql from the browser
func GetGraphiQL(w http.ResponseWriter, r *http.Request) {
	var publicPath = config.PATH + "/graphiql"
	var path = "." + r.URL.Path
	if length := len(path); path[length-1] == '/' {
		path = path[0 : length-1]
	}

	if _, err := os.Stat(path); err != nil {
		http.ServeFile(w, r, publicPath)
		return
	}
	http.ServeFile(w, r, path)
}

//GetGraphQL Get method to describe and update data
func GetGraphQL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()["query"]
	helper.Cors(w, r)
	if len(query) == 0 {
		fmt.Fprint(w, "{\"success\": false, \"error\": \"Mal formado\"}")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result := graph.ExecuteQuery(query[0])
	json.NewEncoder(w).Encode(result)
}

//PostGraphQL Get method to describe and update data
func PostGraphQL(w http.ResponseWriter, r *http.Request) {
	params := helper.GetPostParams(r)
	query := params.Get("query")
	helper.Cors(w, r)
	w.Header().Set("Content-Type", "application/json")
	if query == "" {
		fmt.Fprint(w, "{\"success\": false, \"error\": \"Mal formado\"}")
		return
	}
	defer r.Body.Close()
	result := graph.ExecuteQuery(query)
	if result.HasErrors() {
		fmt.Println(result.Errors)
	}
	json.NewEncoder(w).Encode(result)
}
