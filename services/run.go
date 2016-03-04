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

package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Nothing() {
	log.Printf("Reached services.Nothing endpoint.\n")
}

func Run(binary []byte, tests []string) ([]string, []RunDetails) {
	log.Printf("Reached Run endpoint.\n")

	origDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	tmpDir := os.TempDir()

	if err = os.Chdir(tmpDir); err != nil {
		log.Println(err)
	}

	fileName := fmt.Sprintf("exec%d", time.Now().Nanosecond)
	if err = ioutil.WriteFile(filepath.Join(tmpDir, fileName), binary, 0700); err != nil {
		log.Println(err)
	}

	log.Printf("Judging program.\n")
	out := make([]string, len(tests))
	runinfo := make([]RunDetails, len(tests))
	for i, test := range tests {
		out[i], runinfo[i] = judge(tmpDir, fileName, test)
	}

	if err = os.Chdir(origDir); err != nil {
		log.Println(err)
	}
	return out, runinfo
}
