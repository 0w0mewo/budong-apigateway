package main

import (
	"os"

	"github.com/0w0mewo/budong-apigateway/server/httpserver"
	"github.com/0w0mewo/budong-apigateway/utils"
)

func main() {
	addr := os.Getenv("LOCAL_ADDR")
	if addr == "" {
		addr = ":9988"
	}

	caddr := os.Getenv("SETU_ADDR")
	if caddr == "" {
		caddr = "127.0.0.1:9999"
	}

	apiserver := httpserver.NewRestfulServer(addr, caddr)
	apiserver.Init()
	go apiserver.Run()

	<-utils.WaitForSignal()

	apiserver.Shutdown()

}
