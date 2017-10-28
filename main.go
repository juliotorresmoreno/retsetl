package main

import "bitbucket.org/mlsdatatools/retsetl/bootstrap"
import "bitbucket.org/mlsdatatools/retsetl/taskmgr"
import "bitbucket.org/mlsdatatools/retsetl/ws"

//This is the beginning of the application
func main() {
	taskmgr.Start()
	ws.GetHub().Start()
	bootstrap.StartHTTP()
}
