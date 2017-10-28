package bootstrap

import (
	"net/http"
	"time"

	"bitbucket.org/mlsdatatools/retsetl/config"
	"bitbucket.org/mlsdatatools/retsetl/routes"
)

// StartHTTP Here is defined the start of the application, in this place can make system preconfigurations,
// as well as establish different routers for both http and console responses
func StartHTTP() {
	var mux = routes.GetRoutes()
	var addr = ":" + config.PORT
	var server = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    config.READ_TIMEOUT * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	println("Listening on " + addr)
	println(server.ListenAndServe())
}
