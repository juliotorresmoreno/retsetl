package controllers

import (
	"net/http"

	"bitbucket.org/mlsdatatools/retsetl/ws"
)

//ServerWS Establishes a router capable of handling connections under the websocket protocol
func ServerWS(w http.ResponseWriter, r *http.Request) {
	hub := ws.GetHub()
	hub.ServeWs(w, r)
}
