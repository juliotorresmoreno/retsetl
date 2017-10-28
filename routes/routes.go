package routes

import (
	"net/http"

	"bitbucket.org/mlsdatatools/retsetl/controllers"
	"bitbucket.org/mlsdatatools/retsetl/ws"
	"github.com/gorilla/mux"
)

// GetRoutes Gets the system router
// In this section should be defined each and every one of the anchor points of the application
// This function returns a serverMux object
func GetRoutes() *mux.Router {
	hub := ws.GetHub()
	ServerMux := mux.NewRouter().StrictSlash(false)
	ServerMux.PathPrefix("/graphiql").HandlerFunc(controllers.GetGraphiQL).Methods("GET")
	ServerMux.HandleFunc("/graphql", controllers.GetGraphQL).Methods("GET")
	ServerMux.HandleFunc("/graphql", controllers.PostGraphQL).Methods("POST")
	ServerMux.HandleFunc("/api/v1/data", controllers.GetData).Methods("GET")
	ServerMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r)
	}).Methods("GET")
	ServerMux.PathPrefix("/").HandlerFunc(controllers.GetIndex).Methods("GET")
	return ServerMux
}
