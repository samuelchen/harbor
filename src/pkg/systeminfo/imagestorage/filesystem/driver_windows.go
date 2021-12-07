// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build windows

package filesystem

import (
	"syscall"
	"unsafe"

	"github.com/goharbor/harbor/src/lib/log"
	storage "github.com/goharbor/harbor/src/pkg/systeminfo/imagestorage"
)

const (
	driverName = "filesystem"
)

type driver struct {
	path string
}

// NewDriver returns an instance of filesystem driver
func NewDriver(path string) storage.Driver {
	return &driver{
		path: path,
	}
}

// Name returns a human-readable name of the fielsystem driver
func (d *driver) Name() string {
	return driverName
}

func (d *driver) Cap() (*storage.Capacity, error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		log.Warningf("The path %s is not found, will return zero value of capacity", d.path)
		return &storage.Capacity{Total: 0, Free: 0}, nil
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	if err != nil {
		log.Warningf("GetDiskFreeSpaceExW function not found")
		return &storage.Capacity{Total: 0, Free: 0}, nil
	}

	lpFreeBytesAvailable := uint64(0)
	lpTotalNumberOfBytes := uint64(0)
	lpTotalNumberOfFreeBytes := uint64(0)
	r, a, b := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(d.path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	log.Debugf("r, a, b: %d %d %d", r, a, b)
	// log.Printf("Available  %dmb", lpFreeBytesAvailable/1024/1024.0)
	// log.Printf("Total      %dmb", lpTotalNumberOfBytes/1024/1024.0)
	// log.Printf("Free       %dmb", lpTotalNumberOfFreeBytes/1024/1024.0)

	return &storage.Capacity{
		Total: lpTotalNumberOfBytes,
		Free:  lpTotalNumberOfFreeBytes,
	}, nil
}
