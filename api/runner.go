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
