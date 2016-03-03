package main

import (
	"io/ioutil"
	"log"
	"net/rpc"

	"github.com/clashr/judgrpcd/api"
)

func main() {
	//make connection to rpc server
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatalf("Error in dialing. %s", err)
	}
	//make arguments object
	bin, err := ioutil.ReadFile("a.out")

	tests := make([]api.Test, 3)
	tests[0] = api.Test{"asdf", ""}

	args := &api.Args{"c", bin, tests}
	//this will store returned result
	var result api.Result
	//call remote procedure with args
	if err = client.Call("Judge.Runner", args, &result); err != nil {
		log.Fatalf("Error in running: %s", err)
	}
	//Print out the result
	log.Printf("Mem Used: %dKB\nTime Used: %s\n", result.Data[0].TotalMem,
		result.Data[0].TotalTime)
}
