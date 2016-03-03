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
