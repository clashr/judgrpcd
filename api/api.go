package api

import "time"

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
	TotalMem  int64
}

type Test struct {
	In  string
	Out string
}

func (j *Judge) Interpret(args Args, result *Result) error {
	return interpret(args, result)
}
func (j *Judge) Runner(args Args, result *Result) error {
	return runner(args, result)
}
