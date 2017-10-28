package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"bitbucket.org/mlsdatatools/retsetl/config"

	"bitbucket.org/mlsdatatools/retsetl/helper"
)

//GetIndex Get the start page of the application
func GetIndex(w http.ResponseWriter, r *http.Request) {
	isRets := r.URL.Query()["is_rets"]
	if len(isRets) > 0 {
		helper.Cors(w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
		return
	}
	publicPath := config.PATH + "/public"
	path := publicPath + r.URL.Path
	if _, err := os.Stat(path); err != nil {
		http.ServeFile(w, r, publicPath)
		return
	}
	http.ServeFile(w, r, path)
}
