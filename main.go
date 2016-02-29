package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/clashr/judgrpcd/api"
)

func init() {
}

func main() {
	judge := new(api.Judge)
	err := rpc.Register(judge)
	if err != nil {
		log.Fatalf("Format of service builder isn't correct. %s", err)
	}
	rpc.HandleHTTP()
	//start listening for messages on port 1234
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Couldn't start listening on port 1234. Error %s", err)
	}
	log.Println("Serving RPC handler")
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}
}
