package api

import (
	"log"

	"github.com/clashr/judgrpcd/services"
)

func runner(args Args, result *Result) error {
	log.Printf("Reached Runner Endpoint\n")
	var input []string
	for i, test := range args.TestData {
		input[i] = test.In
	}

	output, runDetails := services.Run(args.Binary, input)

	var sysinfo []Usage
	for i, run := range runDetails {
		sysinfo[i].TotalTime = run.TTotal
		sysinfo[i].TotalMem = run.Mem
	}

	var pass []bool
	for i, test := range args.TestData {
		pass[i] = false
		if output[i] == test.Out {
			pass[i] = true
		}
	}
	*result = Result{pass, 0, sysinfo}
	return nil
}
