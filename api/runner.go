package api

import (
	"log"

	"github.com/clashr/judgrpcd/services"
)

func runner(args Args, result *Result) error {
	log.Printf("Reached Runner Endpoint\n")
	input := make([]string, len(args.TestData))
	for i, test := range args.TestData {
		input[i] = test.In
	}

	output, runDetails := services.Run(args.Binary, input)

	sysinfo := make([]Usage, len(args.TestData))
	for i, run := range runDetails {
		sysinfo[i].TotalTime = run.TTotal
		sysinfo[i].TotalMem = run.Mem
	}

	pass := make([]bool, len(args.TestData))
	for i, test := range args.TestData {
		pass[i] = false
		if output[i] == test.Out {
			pass[i] = true
		}
	}
	*result = Result{pass, 0, sysinfo}
	return nil
}
