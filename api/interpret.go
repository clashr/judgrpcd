package api

import "log"

func interpret(args Args, result *Result) error {
	log.Printf("Reached Interpreter Endpoint\n")
	*result = Result{nil, 0, nil}
	return nil
}
