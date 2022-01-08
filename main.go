package main

import (
	"github.com/0w0mewo/budong-apigateway/server/httpserver"
	"github.com/0w0mewo/budong-apigateway/utils"
)

func main() {
	apiserver := httpserver.NewRestfulServer(":9988")
	apiserver.Init()
	go apiserver.Run()

	<-utils.WaitForSignal()

	apiserver.Shutdown()

}
