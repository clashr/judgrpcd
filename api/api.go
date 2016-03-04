//
// Copyright (c) 2016 Dennis Chen
//
// This file is part of Clashr.
//
// Clashr is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// Clashr is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for
// more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Clashr.  If not, see <http://www.gnu.org/licenses/>.
//

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
