package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	//"github.com/vbatts/go-cgroup"
)

type RunDetails struct {
	TUser  time.Duration
	TSys   time.Duration
	TTotal time.Duration
	Mem    int64
	Exit   int
}

func judge(wd, executable string, testcase string) (string, RunDetails) {
	//TODO: Experiment with CGROUPS
	//cgroup.Init()
	//ns := cgroup.NewCgroup(fmt.Sprintf("/judgrpcd/%d", exec))
	//ns.Create()

	log.Println("Writing Test File to Disk.")
	ioutil.WriteFile(filepath.Join(wd, "data.txt"), []byte(testcase), 0600)
	env := backupenv()
	os.Clearenv()

	log.Println("Starting Executable.")
	elf := fmt.Sprintf("./%s", executable)
	cmd := exec.Command(elf)
	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	var rlimit syscall.Rlimit
	rlimit.Cur = 1024 * 10000 // Provide 10MB of RAM
	rlimit.Max = 1024 * 10001
	log.Println("Limiting Memory to 10MB.")
	log.Printf("New process %s\n", cmd.Process)
	log.Printf("New process pid %d\n", cmd.Process.Pid)
	limMem(cmd.Process.Pid, &rlimit)

	var details RunDetails
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	log.Println("Executable Running... ")
	select {
	case <-time.After(5 * time.Minute): // Limit runtime to 5 minutes
		if err := cmd.Process.Kill(); err != nil {
			log.Println("Failed to Kill: ", err)
		}
		log.Println("Running Program Killed as Timeout Reached")
		details.Exit = -1
	case err := <-done:
		if err != nil {
			log.Printf("Process done with error = %v", err)
			if exiterr, ok := err.(*exec.ExitError); ok {
				// The program has exited with an exit code != 0
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					details.Exit = status.ExitStatus()
				}
			} else {
				details.Exit = -1
			}
		} else {
			details.Exit = 0
		}
	}
	log.Println("Successful Run")

	details.TSys = cmd.ProcessState.SystemTime()
	details.TUser = cmd.ProcessState.UserTime()
	details.TTotal = details.TSys + details.TUser

	details.Mem = cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss

	readBuffer := make([]byte, 1000)
	outFile, err := os.Open(filepath.Join(wd, "output.txt"))
	if err != nil {
		log.Println(err)
	}
	_, err = outFile.Read(readBuffer)
	if err != io.EOF {
		log.Println(err)
	}

	restoreenv(env)

	return string(readBuffer), details
}

func backupenv() map[string]string {
	env := os.Environ()
	pairs := make(map[string]string)
	for _, pair := range env {
		keyValue := strings.Split(pair, "=")
		pairs[keyValue[0]] = keyValue[1]
	}
	return pairs
}

func restoreenv(pairs map[string]string) {
	for key, value := range pairs {
		if err := os.Setenv(key, value); err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func limMem(pid int, rlimit *syscall.Rlimit) error {
	log.Println("Unsafe memory limiting.")
	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRLIMIT64, uintptr(pid),
		syscall.RLIMIT_AS, uintptr(unsafe.Pointer(rlimit)), 0, 0, 0)
	log.Println("Unsafe memory limiting success.")
	var err error
	if errno != 0 {
		err = errno
		return err
	}
	return nil
}
