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

// Package cgroup provides a large amount of functionality for properly
// implementing cgroups in the judging and tracking of programs in Clashr.
// Most of this work is inspired by the DOMJudge implementation runguard.c and
// attempts to mimic a majoritiy of it's behavior
//
package cgroup

// #cgo pkg-config: libcgroup
// #include <libcgroup.h>
import "C"

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

// A uniform interface for the cgroup actions.
type Cgroup struct {
	name string
	cpus string
	mems int64
}

// Init must be called for the cgroup library to be avaliable.
func Init() error {
	if ret := C.cgroup_init(); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Could not cgroup_init %s(%d)", desc, ret)
		return err
	}
	return nil
}

// OutputStats returns the memory used (in bytes) by a process controlled by a
// cgroup.
func (g *Cgroup) OutputStats() (int64, error) {
	var param *C.char

	param = C.CString(g.name)
	cg := C.cgroup_new_cgroup(param)
	if cg == nil {
		return 0, errors.New("cgroup_new_cgroup")
	}
	defer C.cgroup_free(&cg)

	if ret := C.cgroup_get_cgroup(cg); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Get cgroup information: %s(d)", desc, ret)
		return 0, err
	}
	param = C.CString("memory")
	cgController := C.cgroup_get_controller(cg, param)

	var maxUsage C.int64_t
	param = C.CString("memory.memsw.max_usage_in_bytes")
	if ret := C.cgroup_get_value_int64(cgController, param, &maxUsage); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Get cgroup value: %s(%d)", desc, ret)
		return 0, err
	}

	return int64(maxUsage), nil
}

// Create makes a new instance of the cgroup interface. It sets the memory
// limit of the cgroup to memsize, and the set of CPUs to use to cpuset. If
// memory does not need to be limited, set memsize to syscall.RLIMIT_INFINITY.
// Similarly, set cpuset to the empty string if CPUs do not need to be set.
func Create(cgroupname, cpuset string, memsize int64) (Cgroup, error) {
	var param *C.char

	param = C.CString(cgroupname)
	cg := C.cgroup_new_cgroup(param)
	if cg == nil {
		return Cgroup{}, errors.New("cgroup_new_cgroup")
	}
	defer C.cgroup_free(&cg)

	// Set memory restrictions. We limit both ram and ram+swap to the same
	// amount to prevent swapping.
	param = C.CString("memory")
	cgController := C.cgroup_add_controller(cg, param)

	mem := C.int64_t(memsize)
	param = C.CString("memory.limit_in_bytes")
	C.cgroup_add_value_int64(cgController, param, mem)
	param = C.CString("memory.memsw.limit_in_bytes")
	C.cgroup_add_value_int64(cgController, param, mem)

	// Setup CPU restrictions. The task is pinned to a specific set of CPUs
	// if a non-empty string is passed.
	if len(cpuset) > 0 {
		param = C.CString("cpuset")
		cgController = C.cgroup_add_controller(cg, param)

		var arg *C.char
		param = C.CString("cpuset.mems")
		arg = C.CString("0")
		C.cgroup_add_value_string(cgController, param, arg)
		param = C.CString("cpuset.cpus")
		arg = C.CString(cpuset)
		C.cgroup_add_value_string(cgController, param, arg)
	}

	// Perform the creation of the cgroup.
	if ret := C.cgroup_create_cgroup(cg, 1); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Creating cgroup: %s(%d)", desc, ret)
		return Cgroup{}, err
	}
	return Cgroup{cgroupname, cpuset, memsize}, nil
}

// Attach attaches a task with PID p to the cgroup.
func (g *Cgroup) Attach(p int) error {
	var param *C.char

	param = C.CString(g.name)
	cg := C.cgroup_new_cgroup(param)
	if cg == nil {
		return errors.New("cgroup_new_cgroup")
	}
	defer C.cgroup_free(&cg)

	if ret := C.cgroup_get_cgroup(cg); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Get cgroup information: %s(d)", desc, ret)
		return err
	}

	// Attach task to the cgroup
	if ret := C.cgroup_attach_task_pid(cg, C.pid_t(p)); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Attach task to cgroup: %s(%d)", desc, ret)
		return err
	}

	return nil
}

// Kill sends SIGKILL(9) to all PIDs associated with the given cgroup.
func (g *Cgroup) Kill() error {
	cgName := C.CString(g.name)
	cgCtrl := C.CString("memory")

	var pid C.pid_t
	var handle unsafe.Pointer

	for {
		ret := C.cgroup_get_task_begin(cgName, cgCtrl, &handle, &pid)
		C.cgroup_get_task_end(&handle)
		if ret == C.ECGEOF {
			break
		}
		if err := syscall.Kill(int(pid), syscall.SIGKILL); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes the cgroup from the kernel.
func (g *Cgroup) Delete() error {
	var param *C.char

	param = C.CString(g.name)
	cg := C.cgroup_new_cgroup(param)
	if cg == nil {
		err := fmt.Errorf("cgroup_new_cgroup")
		return err
	}

	if ret := C.cgroup_get_cgroup(cg); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Get cgroup information: %s(%d)", desc, ret)
		return err
	}

	// Clean up the cgroup
	if ret := C.cgroup_delete_cgroup(cg, 1); ret != 0 {
		desc := C.GoString(C.cgroup_strerror(ret))
		err := fmt.Errorf("Deleteing cgroup: %s(%d)", desc, ret)
		return err
	}
	return nil
}
