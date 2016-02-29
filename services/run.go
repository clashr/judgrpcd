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

func Run(binary []byte, tests []string) (out []string, err error) {
	log.Printf("Reached Run endpoint.\n")

	origDir, err := os.Getwd()
	if err != nil {
		return
	}

	tmpDir := os.TempDir()

	if err = os.Chdir(tmpDir); err != nil {
		log.Print(err)
		return
	}

	fileName := fmt.Sprintf("exec%d", time.Now().Nanosecond)
	if err = ioutil.WriteFile(filepath.Join(tmpDir, fileName), binary, 0700); err != nil {
		log.Print(err)
		return
	}

	log.Printf("Judging program.\n")
	for i, test := range tests {
		out[i], _ = judge(tmpDir, fileName, test)
	}

	if err = os.Chdir(origDir); err != nil {
		log.Print(err)
		return
	}
	return
}
