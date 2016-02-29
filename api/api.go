package api

import (
	"log"
	"time"

	"github.com/clashr/judgrpcd/services"
)

type Args struct {
	Language string
	Binary   []byte
	TestData []Test
}

type Judge struct {
}

type Result struct {
	Pass  []bool
	Score int
	Data  []Usage
}

//type Cgroups struct {
//	TODO: Complete CGROUPS information
//}

type Usage struct {
	TotalTime time.Duration
	TotalMem  int
}

type Test struct {
	In  string
	Out string
}

func (j *Judge) Interpret(args Args, result *Result) error {
	return Interpret(args, result)
}
func (j *Judge) Runner(args Args, result *Result) error {
	return Runner(args, result)
}

func Interpret(args Args, result *Result) error {
	log.Printf("Reached Interpreter Endpoint\n")
	*result = Result{nil, 0, Cgroups{}, 0}
	return nil
}

func Runner(args Args, result *Result) error {
	log.Printf("Reached Runner Endpoint\n")
	var input []string
	for i, test := range args.TestData {
		input[i] = test.In
	}
	services.Run(args.Binary, input)
	*result = Result{nil, 0, Cgroups{}, 0}
	return nil
}
